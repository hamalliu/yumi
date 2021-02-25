package resolver

import (
	"context"
	"sync"

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
	ss := int64(50)
	r := &Resolver{
		nr:   b.Builder.Build(target, internal.Subset(int(ss))),
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
		if inss, ok := r.nr.Fetch(context.Background()); ok {
			addrs := internal.ToGrpcAddress(inss)
			r.cc.UpdateState(resolver.State{Addresses: addrs})
		}
	}
}
