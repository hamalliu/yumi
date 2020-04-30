package httpsrv

import (
	"net/http"
	"regexp"
	"sync/atomic"
	"time"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// Handler responds to an HTTP request.
type Handler interface {
	ServeHTTP(c *Context)
}

// HandlerFunc http request handler function.
type HandlerFunc func(*Context)

// ServeHTTP calls f(ctx).
func (f HandlerFunc) ServeHTTP(c *Context) {
	f(c)
}

type injection struct {
	pattern  *regexp.Regexp
	handlers []HandlerFunc
}

type Server struct {
	mux *Mux

	server atomic.Value

	Addr         string
	Timeout      time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewServer() *Server {
	//TODO
	return nil
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

func (srv *Server) Run() error {
	//TODO
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
