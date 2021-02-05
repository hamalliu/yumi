package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"yumi/pkg/contextx"
)

const defaultTimeout = time.Second * 2

// TimeoutInterceptor ...
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := contextx.ShrinkDeadline(ctx, timeout)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
