package grpcwrapper

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"yumi/grpcwrapper/resolver/discovery"
)

// clientOptions is gRPC Client
type clientOptions struct {
	discovery     discovery.Builder
	discoveryOpts []discovery.BuildOption

	timeout time.Duration

	grpcOpts []grpc.DialOption
}

// ClientOption is wrapper client option.
type ClientOption func(o *clientOptions)

// WithDiscovery is registed a resolver of grpc
func WithDiscovery(discovery discovery.Builder, opts ...discovery.BuildOption) ClientOption {
	return func(o *clientOptions) {
		o.discovery = discovery
		o.discoveryOpts = opts
	}
}

// WithTimeout with client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithOptions with gRPC options.
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

// Dial returns a GRPC connection.
func Dial(ctx context.Context, target string, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, target, opts...)
}

// DialInsecure returns an insecure GRPC connection.
func DialInsecure(ctx context.Context, target string, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, target, opts...)
}

func dial(ctx context.Context, insecure bool, target string, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout: 500 * time.Millisecond,
	}
	for _, o := range opts {
		o(&options)
	}
	var grpcOpts = []grpc.DialOption{
		grpc.WithTimeout(options.timeout),
	}
	if options.discovery != nil {
		grpc.WithResolvers(discovery.New(options.discovery, options.discoveryOpts...))
	}
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	return grpc.DialContext(ctx, target, grpcOpts...)
}
