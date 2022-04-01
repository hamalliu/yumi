package grpcwrapper

import (
	"context"
	"net"
	"net/url"

	"google.golang.org/grpc"

	"yumi/pkg/log"
)

const loggerName = "transport/grpc"

// ServerOption is gRPC server option.
type ServerOption func(o *Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Options with grpc options.
func Options(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

// Server is a gRPC server wrapper.
type Server struct {
	baseCtx context.Context
	*grpc.Server
	err      error
	lis      net.Listener
	endpoint *url.URL
	network  string
	address  string

	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor
	grpcOpts   []grpc.ServerOption
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
	}
	for _, o := range opts {
		o(srv)
	}
	unaryInts := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	streamInts := []grpc.StreamServerInterceptor{
		srv.streamServerInterceptor(),
	}
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}
	if len(srv.streamInts) > 0 {
		streamInts = append(streamInts, srv.streamInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
		grpc.ChainStreamInterceptor(streamInts...),
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	return srv
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (s *Server) Endpoint() (*url.URL, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.endpoint, nil
}

// Start start the gRPC server.
func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	log.Info("[gRPC] server listening on: ", lis.Addr().String())
	return s.Serve(lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop() error {
	s.GracefulStop()
	log.Info("[gRPC] server stopping")
	return nil
}

// unaryServerInterceptor is a gRPC unary server interceptor
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// ctx, cancel := ic.Merge(ctx, s.baseCtx)
		// defer cancel()
		// md, _ := grpcmd.FromIncomingContext(ctx)
		// replyHeader := grpcmd.MD{}
		// ctx = transport.NewServerContext(ctx, &Transport{
		// 	endpoint:    s.endpoint.String(),
		// 	operation:   info.FullMethod,
		// 	reqHeader:   headerCarrier(md),
		// 	replyHeader: headerCarrier(replyHeader),
		// })
		// if s.timeout > 0 {
		// 	ctx, cancel = context.WithTimeout(ctx, s.timeout)
		// 	defer cancel()
		// }
		// h := func(ctx context.Context, req interface{}) (interface{}, error) {
		// 	return handler(ctx, req)
		// }
		// if len(s.middleware) > 0 {
		// 	h = middleware.Chain(s.middleware...)(h)
		// }
		// reply, err := h(ctx, req)
		// if len(replyHeader) > 0 {
		// 	_ = grpc.SetHeader(ctx, replyHeader)
		// }
		// return reply, err
		return nil, nil
	}
}

// wrappedStream is rewrite grpc stream's context
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func NewWrappedStream(ctx context.Context, stream grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{
		ServerStream: stream,
		ctx:          ctx,
	}
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// streamServerInterceptor is a gRPC stream server interceptor
func (s *Server) streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// ctx, cancel := ic.Merge(ss.Context(), s.baseCtx)
		// defer cancel()
		// md, _ := grpcmd.FromIncomingContext(ctx)
		// replyHeader := grpcmd.MD{}
		// ctx = transport.NewServerContext(ctx, &Transport{
		// 	endpoint:    s.endpoint.String(),
		// 	operation:   info.FullMethod,
		// 	reqHeader:   headerCarrier(md),
		// 	replyHeader: headerCarrier(replyHeader),
		// })

		// ws := NewWrappedStream(ctx, ss)

		// err := handler(srv, ws)
		// if len(replyHeader) > 0 {
		// 	_ = grpc.SetHeader(ctx, replyHeader)
		// }
		// return err
		return nil
	}
}
