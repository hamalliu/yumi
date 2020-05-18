package role

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	rr := r.Group("role")

	rr.Handle(http.MethodPost, "", "add", Add)
	rr.Handle(http.MethodPost, "", "delete", Delete)
	rr.Handle(http.MethodPost, "", "update", Update)
	rr.Handle(http.MethodGet, "", "getitem", GetItem)
	rr.Handle(http.MethodGet, "", "search", Search)
	rr.Handle(http.MethodGet, "", "savepoweruserofrole", SavePowerUserOfRole)
	rr.Handle(http.MethodGet, "", "addpoweruserofrole", AddPowerUserOfRole)
	rr.Handle(http.MethodGet, "", "getpoweruserofrole", GetPowerUserOfRole)
}
