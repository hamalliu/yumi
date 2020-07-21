package gin

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"

	"yumi/pkg/conf"
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
}

func NewServer() *Server {
	srv := &Server{
		Mux: *NewMux(),
	}

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

func (srv *Server) Run(srvConf conf.Server) error {
	server := &http.Server{
		Handler:      srv,
		Addr:         srvConf.Addr,
		ReadTimeout:  srvConf.ReadTimeout.Duration(),
		WriteTimeout: srvConf.WriteTimeout.Duration(),
	}
	srv.server.Store(server)
	_ = server.ListenAndServe()
	return nil
}

func (srv *Server) RunTLS(srvConf conf.Server) error {
	//TODO
	return nil
}

func (srv *Server) RunUnix(srvConf conf.Server) error {
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
