package gin

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"

	"yumi/pkg/conf"
)

//Handler 处理器
type Handler interface {
	ServeHTTP(c *Context)
}

//HandlerFunc 处理函数
type HandlerFunc func(*Context)

func (f HandlerFunc) ServeHTTP(c *Context) {
	f(c)
}

//Server gin服务器
type Server struct {
	Mux

	server atomic.Value
}

//NewServer 新建gin服务器
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

//Run 运行服务器
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

//RunTLS 运行https服务器
func (srv *Server) RunTLS(srvConf conf.Server) error {
	//TODO
	return nil
}

//RunUnix 运行unix域服务器
func (srv *Server) RunUnix(srvConf conf.Server) error {
	//TODO
	return nil
}

//Server ...
func (srv *Server) Server() *http.Server {
	s, ok := srv.server.Load().(*http.Server)
	if !ok {
		return nil
	}
	return s
}

//Shutdown 等待当前处理函数执行完之后，关闭服务器
func (srv *Server) Shutdown(ctx context.Context) error {
	server := srv.Server()
	if server == nil {
		return errors.New("no server")
	}
	return server.Shutdown(ctx)
}
