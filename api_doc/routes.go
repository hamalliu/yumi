package api_doc

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	r.Handle(http.MethodGet, "", "api_doc/*path", ApiDoc)
}