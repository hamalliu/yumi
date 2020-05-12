package menu

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	mr := r.Group("menu")

	mr.Handle(http.MethodPost, "", "add", add)
	mr.Handle(http.MethodPost, "", "delete", delete)
	mr.Handle(http.MethodPut, "", "update", update)
	mr.Handle(http.MethodGet, "", "get", get)
	mr.Handle(http.MethodGet, "", "search", search)
}
