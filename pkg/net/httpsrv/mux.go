package httpsrv

import (
	"net/http"
	"regexp"
)

type Mux struct {
	RouterGroup

	trees methodTrees

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

func NewMux() *Mux {
	mux := &Mux{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		trees:                  make(methodTrees, 0, 9),
		HandleMethodNotAllowed: true,
		injections:             make([]injection, 0),
	}

	mux.RouterGroup.mux = mux

	mux.NoRoute(default404Handler)
	mux.NoMethod(default405Handler)

	return mux
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
		handlers, params, _ := root.getValue(rPath, c.Params, unescape)
		if handlers != nil {
			c.handlers = handlers
			c.Params = params
			return
		}
		break
	}

	if mux.HandleMethodNotAllowed {
		for _, tree := range mux.trees {
			if tree.method == httpMethod {
				continue
			}
			if handlers, _, _ := tree.root.getValue(rPath, nil, unescape); handlers != nil {
				c.handlers = mux.allNoMethod
				return
			}
		}
	}
	c.handlers = mux.allNoRoute
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
	//if _, ok := mux.metastore[path]; !ok {
	//	mux.metastore[path] = make(map[string]interface{})
	//}
	//mux.metastore[path]["method"] = method
	root := mux.trees.get(method)
	if root == nil {
		root = new(node)
		mux.trees = append(mux.trees, methodTree{method: method, root: root})
	}

	//prelude := func(c *Context) {
	//	c.method = method
	//	c.RoutePath = path
	//}
	//handlers = append([]HandlerFunc{prelude}, handlers...)
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
