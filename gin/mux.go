package gin

import (
	"net/http"
	"path"
	"regexp"
)

type injection struct {
	pattern  *regexp.Regexp
	handlers []HandlerFunc
}

//Mux 路由选择器
type Mux struct {
	RouterGroup

	trees methodTrees

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the url.RawPath will be used to find parameters.
	UseRawPath bool

	// If true, the path value will be unescaped.
	// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
	// as url.Path gonna be used, which is already unescaped.
	UnescapePathValues bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	injections []injection

	allNoRoute  []HandlerFunc
	allNoMethod []HandlerFunc
	noRoute     []HandlerFunc
	noMethod    []HandlerFunc
}

func default404Handler(c *Context) {
	c.Bytes(404, "text/plain", []byte(http.StatusText(404)))
	c.Abort()
}

func default405Handler(c *Context) {
	c.Bytes(405, "text/plain", []byte(http.StatusText(405)))
	c.Abort()
}

//NewMux 新建一个路由选择器
func NewMux() *Mux {
	mux := &Mux{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		trees:                  make(methodTrees, 0, 9),
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		UseRawPath:             false,
		UnescapePathValues:     true,
		HandleMethodNotAllowed: true,
		injections:             make([]injection, 0),
	}

	mux.RouterGroup.mux = mux

	mux.NoRoute(default404Handler)
	mux.NoMethod(default405Handler)

	return mux
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := &Context{
		Context:  nil,
		mux:      mux,
		index:    -1,
		handlers: nil,
		Keys:     nil,
	}

	c.Request = req
	c.Writer = w

	mux.handleContext(c)
}

func (mux *Mux) handleContext(c *Context) {
	mux.handleHTTPRequest(c)
	c.Next()
}

func (mux *Mux) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	if mux.UseRawPath && len(c.Request.URL.EscapedPath()) > 0 {
		rPath = c.Request.URL.EscapedPath()
		unescape = mux.UnescapePathValues
	}
	rPath = cleanPath(rPath)

	// Find root of the tree for the given HTTP method
	t := mux.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		value := root.getValue(rPath, c.Params, unescape)
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath
			c.code = value.code
			c.Params = value.params
			c.Next()
			return
		}

		if httpMethod != "CONNECT" && rPath != "/" {
			if value.tsr && mux.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			if mux.RedirectFixedPath && redirectFixedPath(c, root, mux.RedirectFixedPath) {
				return
			}
		}
		break
	}

	if mux.HandleMethodNotAllowed {
		for _, tree := range mux.trees {
			if tree.method == httpMethod {
				continue
			}
			if value := tree.root.getValue(rPath, nil, unescape); value.handlers != nil {
				c.handlers = mux.allNoMethod
				c.Next()
				return
			}
		}
	}
	c.handlers = mux.allNoRoute
	c.Next()
	return
}

func (mux *Mux) addRoute(method, path string, handlers ...HandlerFunc) {
	if path[0] != '/' {
		panic("path must begin with '/'")
	}
	if method == "" {
		panic("HTTP method can not be empty")
	}
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}
	root := mux.trees.get(method)
	if root == nil {
		root = new(node)
		mux.trees = append(mux.trees, methodTree{method: method, root: root})
	}

	root.addRoute(path, handlers)
}

// UseFunc attachs a global middleware to the router. ie. the middleware attached though UseFunc() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (mux *Mux) UseFunc(middleware ...HandlerFunc) IRoutes {
	mux.RouterGroup.UseFunc(middleware...)
	return mux
}

// Use attachs a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (mux *Mux) Use(middleware ...Handler) IRoutes {
	mux.RouterGroup.Use(middleware...)
	return mux
}

// Inject is
func (mux *Mux) Inject(pattern string, handlers ...HandlerFunc) {
	mux.injections = append(mux.injections, injection{
		pattern:  regexp.MustCompile(pattern),
		handlers: handlers,
	})
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (mux *Mux) NoRoute(handlers ...HandlerFunc) {
	mux.noRoute = handlers
	mux.rebuild404Handlers()
}

// NoMethod sets the handlers called when... TODO.
func (mux *Mux) NoMethod(handlers ...HandlerFunc) {
	mux.noMethod = handlers
	mux.rebuild405Handlers()
}

func (mux *Mux) rebuild404Handlers() {
	mux.allNoRoute = mux.combineHandlers(mux.noRoute)
}

func (mux *Mux) rebuild405Handlers() {
	mux.allNoMethod = mux.combineHandlers(mux.noMethod)
}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	p := req.URL.Path
	if prefix := path.Clean(c.Request.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}
	redirectRequest(c)
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	rPath := req.URL.Path

	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(rPath), trailingSlash); ok {
		req.URL.Path = string(fixedPath)
		redirectRequest(c)
		return true
	}
	return false
}

func redirectRequest(c *Context) {
	req := c.Request
	rURL := req.URL.String()

	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}
	http.Redirect(c.Writer, req, rURL, code)
}
