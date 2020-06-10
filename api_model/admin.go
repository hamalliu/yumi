package api_model

import "yumi/api_doc"

var login = api_doc.Component{
	RequestBody: ReqBodyLogin{},
}

type ReqBodyLogin struct {
}

var logout = api_doc.Component{
	RequestBody: ReqBodyLogout{},
}

type ReqBodyLogout struct {
}
