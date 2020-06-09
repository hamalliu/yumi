package api_model

import "yumi/doc"

var login = doc.Component{
	RequestBody: ReqBodyLogin{},
}

type ReqBodyLogin struct {
}

var logout = doc.Component{
	RequestBody: ReqBodyLogout{},
}

type ReqBodyLogout struct {
}
