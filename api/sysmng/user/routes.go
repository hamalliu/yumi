package user

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	rr := r.Group("user")

	rr.Handle(http.MethodPost, "", "/add", Add)
	rr.Handle(http.MethodPost, "", "/delete", Delete)
	rr.Handle(http.MethodPut, "", "/update", Update)
	rr.Handle(http.MethodGet, "", "/GetItem", GetItem)
	rr.Handle(http.MethodGet, "", "/Search", Search)
	rr.Handle(http.MethodPost, "", "/AddPowerRoleOfUser", AddPowerRoleOfUser)
	rr.Handle(http.MethodPost, "", "/SavePowerRoleOfUser", SavePowerRoleOfUser)
	rr.Handle(http.MethodGet, "", "/GetPowerRoleOfUser", GetPowerRoleOfUser)
	rr.Handle(http.MethodGet, "", "/GetAllPowerOfUser", GetAllPowerOfUser)
}
