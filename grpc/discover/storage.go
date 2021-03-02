package discover

import (
	"context"

	"google.golang.org/grpc/resolver"
)

// Registry Register an instance and renew automatically.
type Registry interface {
	Register(ctx context.Context, ins *Instance) (cancel context.CancelFunc, err error)
	Close() error
}

// Resolver resolve naming service
type Resolver interface {
	Fetch(context.Context) ([]*Instance, bool)
	Watch() <-chan struct{}
	Close() error
}

// Builder resolver builder.
type Builder interface {
	Build(target resolver.Target) Resolver
	Scheme() string
}
