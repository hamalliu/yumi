package gin

//Handler 处理器
type Handler interface {
	ServeHTTP(c *Context)
}

//HandlerFunc 处理函数
type HandlerFunc func(*Context)

func (f HandlerFunc) ServeHTTP(c *Context) {
	f(c)
}
