package ymhttp

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"time"
)

type Handler interface {
	ServeHTTP(c *Context)
}

type HandlerFunc func(*Context)

func (f HandlerFunc) ServeHTTP(c *Context) {
	f(c)
}

type Server struct {
	Mux

	server atomic.Value

	conf Config
}

type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func DefalutServer() *Server {
	srv := &Server{
		Mux: *NewMux(),
		conf: Config{
			Addr:         ":8888",
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
		},
	}
	//srv.Use(middeware.Cors(), middeware.Recovery(), middeware.PrintRequest())
	return srv
}

func NewServer(conf Config) *Server {
	srv := &Server{
		Mux:  *NewMux(),
		conf: conf,
	}

	//srv.Use(middeware.Cors(), middeware.Recovery(), middeware.PrintRequest())
	return srv
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := &Context{
		Context:  nil,
		srv:      srv,
		index:    -1,
		handlers: nil,
		Keys:     nil,
	}

	c.Request = req
	c.Writer = w

	srv.handleContext(c)
}

func (srv *Server) handleContext(c *Context) {
	srv.mux.handleHTTPRequest(c)
	c.Next()
}

func (srv *Server) Run(addr string) error {
	if addr != "" {
		srv.conf.Addr = addr
	}
	server := &http.Server{
		Handler:      srv,
		Addr:         srv.conf.Addr,
		ReadTimeout:  srv.conf.ReadTimeout,
		WriteTimeout: srv.conf.WriteTimeout,
	}
	srv.server.Store(server)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	return nil
}

func (srv *Server) RunTLS() error {
	//TODO
	return nil
}

func (srv *Server) RunUnix() error {
	//TODO
	return nil
}

func (srv *Server) Server() *http.Server {
	s, ok := srv.server.Load().(*http.Server)
	if !ok {
		return nil
	}
	return s
}

func (srv *Server) Shutdown(ctx context.Context) error {
	server := srv.Server()
	if server == nil {
		return errors.New("no server")
	}
	return server.Shutdown(ctx)
}
