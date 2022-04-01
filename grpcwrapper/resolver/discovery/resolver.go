package discovery

import (
	"context"
	"sync"

	"google.golang.org/grpc/resolver"
)

var (
	_  resolver.Resolver = &BaseResolver{}
	_  resolver.Builder  = &BaseBuilder{}
	mu sync.Mutex
)

// Register register resolver builder if nil.
func Register(b Builder, opts ...BuildOption) {
	mu.Lock()
	defer mu.Unlock()
	if resolver.Get(b.Scheme()) == nil {
		resolver.Register(&BaseBuilder{b, opts})
	}
}

// BaseBuilder is also a resolver builder.
// It's build() function always returns itself.
type BaseBuilder struct {
	Builder
	opts []BuildOption
}

// Build returns itself for Resolver, because it's both a builder and a resolver.
func (b *BaseBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// ss := int64(50)
	r := &BaseResolver{
		nr:   b.Builder.Build(target),
		cc:   cc,
		quit: make(chan struct{}, 1),
	}
	for i := range b.opts {
		b.opts[i].Apply(&r.opts)
	}
	go r.updateproc()
	return r, nil
}

// BaseResolver watches for the updates on the specified target.
// Updates include address updates and service config updates.
type BaseResolver struct {
	nr           Resolver
	cc           resolver.ClientConn
	quit         chan struct{}
	interceptors []func([]*Instance) []*Instance
	opts         BuildOptions
}

// Close is a noop for Resolver.
func (r *BaseResolver) Close() {
	select {
	case r.quit <- struct{}{}:
		r.nr.Close()
	default:
	}
}

// ResolveNow is a noop for Resolver.
func (r *BaseResolver) ResolveNow(o resolver.ResolveNowOptions) {
}

func (r *BaseResolver) updateproc() {
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
		if inss, ok := r.nr.Fetch(context.Background()); ok {
			if r.opts.Subset != nil {
				inss = r.opts.Subset(inss, r.opts.SubsetSize)
			}
			
			addrs := ToGrpcAddress(inss)
			r.cc.UpdateState(resolver.State{Addresses: addrs})
		}
	}
}
