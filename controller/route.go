package controller

import (
	"fmt"
)

type Route struct {
	Pattern string
}

func (r Route) Group(path string, middlewareFunc ...interface{}) Route {
	r.Pattern = fmt.Sprintf("%s/%s", r.Pattern, path)

	if middlewareFunc != nil {
		ctlr.srve.Group(r.Pattern, middlewareFunc)
	}

	return r
}

func (r Route) Handle(httpMethod, path string, handler interface{}, patternCode string, conf *HandlerConf) {
	r.Pattern = fmt.Sprintf("%s/%s", r.Pattern, path)

	GetHandlerConfs().add(patternCode, conf)

	ctlr.srve.Handle(httpMethod, r.Pattern, handler, patternCode)
}
