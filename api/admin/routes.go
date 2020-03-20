package admin

import (
	"net/http"

	"yumi/controller"
)

func Mount(r controller.Route) {
	ar := r.Group("admin")

	ar.Handle(http.MethodPost, "login", login, "0011", nil)
	ar.Handle(http.MethodPost, "logout", logout, "0011", nil)
}
