package registry

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"

	"yumi/grpcwrapper/resolver/discovery"
)

var (
	//etcdPrefix is a etcd globe key prefix
	endpoints  string
	etcdPrefix string

	//Time units is second
	registerTTL        = 90
	defaultDialTimeout = 30
)

var (
	//ErrDuplication is a register duplication err
	ErrDuplication = errors.New("etcd: instance duplicate registration")
)

func init() {
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	// env
	fs.StringVar(&endpoints, "etcd.endpoints", os.Getenv("ETCD_ENDPOINTS"), "etcd.endpoints is etcd endpoints. value: 127.0.0.1:2379,127.0.0.2:2379 etc.")
	fs.StringVar(&etcdPrefix, "etcd.prefix", defaultString("ETCD_PREFIX", "kratos_etcd"), "etcd globe key prefix or use ETCD_PREFIX env variable. value etcd_prefix etc.")
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}

var _ discovery.Builder = &EtcdBuilder{}

// EtcdBuilder is a etcd clientv3 EtcdBuilder
type EtcdBuilder struct {
	cli   *clientv3.Client
	mutex sync.RWMutex
	apps  map[string]*appInfo
}
type appInfo struct {
	resolver map[*Resolve]struct{}
	ins      atomic.Value
	e        *EtcdBuilder
	once     sync.Once
}

// Resolve etch resolver.
type Resolve struct {
	key   string
	event chan struct{}
	e     *EtcdBuilder
}

// NewEtcdBuilder is new a etcdbuilder
func NewEtcdBuilder(c *clientv3.Config) (e *EtcdBuilder, err error) {
	if c == nil {
		if endpoints == "" {
			panic(fmt.Errorf("invalid etcd config endpoints:%+v", endpoints))
		}
		c = &clientv3.Config{
			Endpoints:   strings.Split(endpoints, ","),
			DialTimeout: time.Second * time.Duration(defaultDialTimeout),
			// DialOptions: []grpc.DialOption{grpc.WithBlock()},
		}
	}
	cli, err := clientv3.New(*c)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// ctx, cancel := context.WithCancel(context.Background())
	e = &EtcdBuilder{
		cli:  cli,
		apps: map[string]*appInfo{},
	}
	return
}

// Build disovery resovler builder.
func (e *EtcdBuilder) Build(target resolver.Target) discovery.Resolver {
	appID := target.Endpoint
	r := &Resolve{
		key:   appID,
		e:     e,
		event: make(chan struct{}, 1),
	}

	e.mutex.Lock()
	app, ok := e.apps[appID]
	if !ok {
		app = &appInfo{
			resolver: make(map[*Resolve]struct{}),
			e:        e,
		}
		e.apps[appID] = app
	}
	app.resolver[r] = struct{}{}
	e.mutex.Unlock()
	if ok {
		select {
		case r.event <- struct{}{}:
		default:
		}
	}

	app.once.Do(func() {
		go app.watch(appID)
		// log.Info("etcd: AddWatch(%s) already watch(%v)", appID, ok)
	})
	return r
}

// Scheme return etcd's scheme
func (e *EtcdBuilder) Scheme() string {
	return "etcd"
}

func (a *appInfo) watch(appID string) {
	_ = a.fetchstore(appID)
	prefix := fmt.Sprintf("/%s/%s/", etcdPrefix, appID)
	rch := a.e.cli.Watch(context.TODO(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type == clientv3.EventTypePut || ev.Type == clientv3.EventTypeDelete {
				_ = a.fetchstore(appID)
			}
		}
	}
}
func (a *appInfo) fetchstore(appID string) (err error) {
	prefix := fmt.Sprintf("/%s/%s/", etcdPrefix, appID)
	resp, err := a.e.cli.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		return errors.WithStack(fmt.Errorf("etcd: fetch client.Get(%s) error(%+v)", prefix, err))
	}
	ins, err := a.paserIns(resp)
	if err != nil {
		return err
	}
	a.store(ins)
	return nil
}
func (a *appInfo) paserIns(resp *clientv3.GetResponse) (ins []*discovery.Instance, err error) {
	for _, ev := range resp.Kvs {
		in := new(discovery.Instance)

		err := json.Unmarshal(ev.Value, in)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ins = append(ins, in)
	}

	return ins, nil
}
func (a *appInfo) store(ins []*discovery.Instance) {
	a.ins.Store(ins)
	a.e.mutex.RLock()
	for rs := range a.resolver {
		select {
		case rs.event <- struct{}{}:
		default:
		}
	}
	a.e.mutex.RUnlock()
}

// Watch watch instance.
func (r *Resolve) Watch() <-chan struct{} {
	return r.event
}

// Fetch fetch resolver instance.
func (r *Resolve) Fetch(ctx context.Context) (ins []*discovery.Instance, ok bool) {
	r.e.mutex.RLock()
	app, ok := r.e.apps[r.key]
	r.e.mutex.RUnlock()
	if ok {
		ins, ok = app.ins.Load().([]*discovery.Instance)
		return
	}
	return
}

// Close close resolver.
func (r *Resolve) Close() error {
	r.e.mutex.Lock()
	if app, ok := r.e.apps[r.key]; ok && len(app.resolver) != 0 {
		delete(app.resolver, r)
	}
	r.e.mutex.Unlock()
	return nil
}

// =======================EtcdRegistry===============================

var _ discovery.Registry = &EtcdRegistry{}

// EtcdRegistry ...
type EtcdRegistry struct {
	mutex      sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	cli        *clientv3.Client
	registry   map[string]struct{}
}

// NewEtcdRegistry is new a EtcdRegistry
func NewEtcdRegistry(c *clientv3.Config) (e *EtcdRegistry, err error) {
	if c == nil {
		if endpoints == "" {
			panic(fmt.Errorf("invalid etcd config endpoints:%+v", endpoints))
		}
		c = &clientv3.Config{
			Endpoints:   strings.Split(endpoints, ","),
			DialTimeout: time.Second * time.Duration(defaultDialTimeout),
			// DialOptions: []grpc.DialOption{grpc.WithBlock()},
		}
	}
	cli, err := clientv3.New(*c)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	e = &EtcdRegistry{
		cli:        cli,
		ctx:        ctx,
		cancelFunc: cancel,
		registry:   make(map[string]struct{}),
	}
	return
}

// Register is register instance
func (e *EtcdRegistry) Register(ctx context.Context, ins *discovery.Instance) (cancelFunc context.CancelFunc, err error) {
	e.mutex.Lock()
	if _, ok := e.registry[ins.AppID]; ok {
		err = ErrDuplication
	} else {
		e.registry[ins.AppID] = struct{}{}
	}
	e.mutex.Unlock()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	ctx, cancel := context.WithCancel(e.ctx)
	if err = e.register(ctx, ins); err != nil {
		e.mutex.Lock()
		delete(e.registry, ins.AppID)
		e.mutex.Unlock()
		cancel()
		err = errors.WithStack(err)
		return
	}
	ch := make(chan struct{}, 1)
	cancelFunc = context.CancelFunc(func() {
		cancel()
		<-ch
	})

	go func() {

		ticker := time.NewTicker(time.Duration(registerTTL/3) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = e.register(ctx, ins)
			case <-ctx.Done():
				_ = e.unregister(ins)
				ch <- struct{}{}
				return
			}
		}
	}()
	return
}

//注册和续约公用一个操作
func (e *EtcdRegistry) register(ctx context.Context, ins *discovery.Instance) (err error) {
	prefix := e.keyPrefix(ins)
	val, _ := json.Marshal(ins)

	ttlResp, err := e.cli.Grant(context.TODO(), int64(registerTTL))
	if err != nil {
		return errors.WithStack(fmt.Errorf("etcd: register client.Lease.Create (%v) error(%v)", registerTTL, err))
	}
	_, err = e.cli.Put(ctx, prefix, string(val), clientv3.WithLease(clientv3.LeaseID(ttlResp.ID)))
	if err != nil {
		err = fmt.Errorf("etcd: register client.Put(%v) appid(%s) hostname(%s) error(%v)",
			prefix, ins.AppID, ins.Hostname, err)
		return errors.WithStack(err)
	}
	return nil
}
func (e *EtcdRegistry) unregister(ins *discovery.Instance) (err error) {
	prefix := e.keyPrefix(ins)

	if _, err = e.cli.Delete(context.TODO(), prefix); err != nil {
		err = fmt.Errorf("etcd: unregister client.Delete(%v) appid(%s) hostname(%s) error(%v)",
			prefix, ins.AppID, ins.Hostname, err)
		return errors.WithStack(err)
	}
	// log.Info("etcd: unregister client.Delete(%v)  appid(%s) hostname(%s) success",
	// 	prefix, ins.AppID, ins.Hostname)
	return
}
func (e *EtcdRegistry) keyPrefix(ins *discovery.Instance) string {
	return fmt.Sprintf("/%s/%s/%s", etcdPrefix, ins.AppID, ins.Hostname)
}

// Close stop all running process including etcdfetch and register
func (e *EtcdRegistry) Close() error {
	e.cancelFunc()
	return nil
}
