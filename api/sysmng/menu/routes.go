package menu

import (
	"net/http"

	"yumi/controller"
)

func Mount(r controller.Route) {
	mr := r.Group("menu")

	mr.Handle(http.MethodPost, "add", add, "0011", nil)
	mr.Handle(http.MethodPost, "delete", delete, "0011", nil)
	mr.Handle(http.MethodPut, "update", update, "0011", nil)
	mr.Handle(http.MethodGet, "get", get, "0011", nil)
	mr.Handle(http.MethodGet, "search", search, "0011", nil)
}
