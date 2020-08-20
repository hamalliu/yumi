package apimodel

import "yumi/apidoc"

var login = apidoc.Component{
	RequestBody: ReqBodyLogin{},
}

//ReqBodyLogin ...
type ReqBodyLogin struct {
}

var logout = apidoc.Component{
	RequestBody: ReqBodyLogout{},
}

//ReqBodyLogout ...
type ReqBodyLogout struct {
}
