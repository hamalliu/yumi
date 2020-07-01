package api_doc

import (
	"net/http"

	"yumi/pkg/net/gin"
)

func Mount(r *gin.RouterGroup) {
	r.Handle(http.MethodGet, "", "api_doc/*path", ApiDoc)
}
