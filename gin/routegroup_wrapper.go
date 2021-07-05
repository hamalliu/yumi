package gin

type GroupRoutes interface {
	Group(string, string, ...HandlerFunc) *RouterGroup
	
	// Handle(string, string, string, ...HandlerFunc) IRoutes
	HEAD(string, string, ...HandlerFunc) IRoutes
	GET(string, string, ...HandlerFunc) IRoutes
	POST(string, string, ...HandlerFunc) IRoutes
	PUT(string, string, ...HandlerFunc) IRoutes
	DELETE(string, string, ...HandlerFunc) IRoutes
}
