package resolver

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"

	"yumi/grpc/resolver/internal"
)

var (
	_  resolver.Resolver = &Resolver{}
	_  resolver.Builder  = &Builder{}
	mu sync.Mutex
)

// Register register resolver builder if nil.
func Register(b internal.Builder) {
	mu.Lock()
	defer mu.Unlock()
	if resolver.Get(b.Scheme()) == nil {
		resolver.Register(&Builder{b})
	}
}

// Set override any registered builder
func Set(b internal.Builder) {
	mu.Lock()
	defer mu.Unlock()
	resolver.Register(&Builder{b})
}

// Builder is also a resolver builder.
// It's build() function always returns itself.
type Builder struct {
	internal.Builder
}

// Build returns itself for Resolver, because it's both a builder and a resolver.
func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var zone string
	ss := int64(50)
	clusters := map[string]struct{}{}
	str := strings.SplitN(target.Endpoint, "?", 2)
	if len(str) == 0 {
		return nil, errors.Errorf("warden resolver: parse target.Endpoint(%s) failed!err:=endpoint is empty", target.Endpoint)
	} else if len(str) == 2 {
		m, err := url.ParseQuery(str[1])
		if err == nil {
			for _, c := range m[internal.MetaCluster] {
				clusters[c] = struct{}{}
			}
			zones := m[internal.MetaZone]
			if len(zones) > 0 {
				zone = zones[0]
			}
			if sub, ok := m["subset"]; ok {
				if t, err := strconv.ParseInt(sub[0], 10, 64); err == nil {
					ss = t
				}
			}
		}
	}
	r := &Resolver{
		nr:   b.Builder.Build(str[0], internal.Filter(Scheme, clusters), internal.ScheduleNode(zone), internal.Subset(int(ss))),
		cc:   cc,
		quit: make(chan struct{}, 1),
	}
	go r.updateproc()
	return r, nil
}

// Resolver watches for the updates on the specified target.
// Updates include address updates and service config updates.
type Resolver struct {
	nr   internal.Resolver
	cc   resolver.ClientConn
	quit chan struct{}
}

// Close is a noop for Resolver.
func (r *Resolver) Close() {
	select {
	case r.quit <- struct{}{}:
		r.nr.Close()
	default:
	}
}

// ResolveNow is a noop for Resolver.
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (r *Resolver) updateproc() {
	event := r.nr.Watch()
	for {
		select {
		case <-r.quit:
			return
		case _, ok := <-event:
			if !ok {
				return
			}
		}
		if ins, ok := r.nr.Fetch(context.Background()); ok {
			addrs := internal.ToGrpcAddress(ins.Instances)
			r.cc.UpdateState(resolver.State{Addresses: addrs})
		}
	}
}
