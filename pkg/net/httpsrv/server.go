package httpsrv

import (
	"net/http"
	"regexp"
	"sync"
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

// MethodConfig is
type MethodConfig struct {
	Timeout time.Duration
}

type injection struct {
	pattern  *regexp.Regexp
	handlers []HandlerFunc
}

type Server struct {
	mux *Mux

	server atomic.Value

	pcLock        sync.RWMutex
	methodConfigs map[string]*MethodConfig

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
	//var cancel func()
	//req := c.Request
	//ctype := req.Header.Get("Content-Type")
	//switch {
	//case strings.Contains(ctype, "multipart/form-data"):
	//	req.ParseMultipartForm(defaultMaxMemory)
	//default:
	//	req.ParseForm()
	//}
	//// get derived timeout from http request header,
	//// compare with the engine configured,
	//// and use the minimum one
	//srv.lock.RLock()
	//tm := time.Duration(engine.conf.Timeout)
	//srv.lock.RUnlock()
	//// the method config is preferred
	//if pc := srv.methodConfig(c.Request.URL.Path); pc != nil {
	//	tm = time.Duration(pc.Timeout)
	//}
	//if ctm := timeout(req); ctm > 0 && tm > ctm {
	//	tm = ctm
	//}
	////md := metadata.MD{
	////	metadata.RemoteIP:    remoteIP(req),
	////	metadata.RemotePort:  remotePort(req),
	////	metadata.Criticality: string(criticality.Critical),
	////}
	//parseMetadataTo(req, md)
	//ctx := metadata.NewContext(context.Background(), md)
	//if tm > 0 {
	//	c.Context, cancel = context.WithTimeout(ctx, tm)
	//} else {
	//	c.Context, cancel = context.WithCancel(ctx)
	//}
	//defer cancel()
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
