package gin

import (
	"fmt"

	"yumi/pkg/log"
)

// IRouter http router framework interface.
type IRouter interface {
	IRoutes
	Group(string, string, ...HandlerFunc) *RouterGroup
}

// IRoutes http router interface.
type IRoutes interface {
	UseFunc(...HandlerFunc) IRoutes
	Use(...Handler) IRoutes

	// Handle(string, string, string, ...HandlerFunc) IRoutes
	HEAD(string, string, ...HandlerFunc) IRoutes
	GET(string, string, ...HandlerFunc) IRoutes
	POST(string, string, ...HandlerFunc) IRoutes
	PUT(string, string, ...HandlerFunc) IRoutes
	DELETE(string, string, ...HandlerFunc) IRoutes
}

// RouterGroup is used internally to configure router, a RouterGroup is associated with a prefix
// and an array of handlers (middleware).
type RouterGroup struct {
	Handlers    []HandlerFunc
	basePath    string
	description string
	mux         *Mux
	root        bool
}

var _ IRouter = &RouterGroup{}

// Use adds middleware to the group, see example code in api_doc.
func (group *RouterGroup) Use(middleware ...Handler) IRoutes {
	for _, m := range middleware {
		group.Handlers = append(group.Handlers, m.ServeHTTP)
	}
	return group.returnObj()
}

// UseFunc adds middleware to the group, see example code in api_doc.
func (group *RouterGroup) UseFunc(middleware ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnObj()
}

// Group creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (group *RouterGroup) Group(description, relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		description: description,
		basePath: group.calculateAbsolutePath(relativePath),
		mux:      group.mux,
		root:     false,
	}
}

// BasePath router group base path.
func (group *RouterGroup) BasePath() string {
	return group.basePath
}

func (group *RouterGroup) handle(description, httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	injections := group.injections(relativePath)
	handlers = group.combineHandlers(injections, handlers)
	group.mux.addRoute(httpMethod, absolutePath, handlers...)
	group.description = description
	log.Debug(fmt.Sprintf("method:%s, path:%s", httpMethod, absolutePath))
	return group.returnObj()
}

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in api_doc.
//
// For HEAD, GET, POST, PUT, and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
//func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
//	if matches, err := regexp.MatchString("^[A-Z]+$", httpMethod); !matches || err != nil {
//		panic("http method " + httpMethod + " is not valid")
//	}
//	return group.handle(httpMethod, relativePath, handlers...)
//}

// GET is a shortcut for router.Handle("GET", path, handle).
func (group *RouterGroup) GET(description, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(description, "GET", relativePath, handlers...)
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (group *RouterGroup) POST(description, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(description, "POST", relativePath, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (group *RouterGroup) PUT(description, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(description, "PUT", relativePath, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (group *RouterGroup) DELETE(description, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(description, "DELETE", relativePath, handlers...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (group *RouterGroup) PATCH(description, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(description, "PATCH", relativePath, handlers...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (group *RouterGroup) OPTIONS(description, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(description, "OPTIONS", relativePath, handlers...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (group *RouterGroup) HEAD(description, relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(description, "HEAD", relativePath, handlers...)
}

func (group *RouterGroup) combineHandlers(handlerGroups ...[]HandlerFunc) []HandlerFunc {
	finalSize := len(group.Handlers)
	for _, handlers := range handlerGroups {
		finalSize += len(handlers)
	}
	if finalSize >= int(_abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, group.Handlers)
	position := len(group.Handlers)
	for _, handlers := range handlerGroups {
		copy(mergedHandlers[position:], handlers)
		position += len(handlers)
	}
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.mux
	}
	return group
}

// injections is
func (group *RouterGroup) injections(relativePath string) []HandlerFunc {
	absPath := group.calculateAbsolutePath(relativePath)
	for _, injection := range group.mux.injections {
		if !injection.pattern.MatchString(absPath) {
			continue
		}
		return injection.handlers
	}
	return nil
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (group *RouterGroup) Any(description, relativePath string, handlers ...HandlerFunc) IRoutes {
	group.handle(description, "GET", relativePath, handlers...)
	group.handle(description, "POST", relativePath, handlers...)
	group.handle(description, "PUT", relativePath, handlers...)
	group.handle(description, "PATCH", relativePath, handlers...)
	group.handle(description, "HEAD", relativePath, handlers...)
	group.handle(description, "OPTIONS", relativePath, handlers...)
	group.handle(description, "DELETE", relativePath, handlers...)
	group.handle(description, "CONNECT", relativePath, handlers...)
	group.handle(description, "TRACE", relativePath, handlers...)
	return group.returnObj()
}
