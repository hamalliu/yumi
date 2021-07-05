package grpcwrapper

import (
	"fmt"
	"net"

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
	*grpc.Server
	lis      net.Listener
	network  string
	address  string
	grpcOpts []grpc.ServerOption
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
	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(),
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
func (s *Server) Endpoint() (string, error) {
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc://%s", addr), nil
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
