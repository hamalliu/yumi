package internal

import (
	"context"
)

// Registry Register an instance and renew automatically.
type Registry interface {
	Register(ctx context.Context, ins *Instance) (cancel context.CancelFunc, err error)
	Close() error
}

// Resolver resolve naming service
type Resolver interface {
	Fetch(context.Context) (*InstancesInfo, bool)
	Watch() <-chan struct{}
	Close() error
}

// Builder resolver builder.
type Builder interface {
	Build(id string, options ...BuildOption) Resolver
	Scheme() string
}
