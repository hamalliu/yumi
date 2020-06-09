package admin

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	ar := r.Group("admin")

	ar.Handle(http.MethodPost, "", "login", login)
	ar.Handle(http.MethodPost, "", "logout", logout)
}
