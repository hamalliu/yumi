package media

import (
	"net/http"
	"yumi/controller"
)

func Mount(r controller.Route) {
	mr := r.Group("media")

	mr.Handle(http.MethodPost, "uploadmultipart", UploadMultipart, "0011", nil)
	mr.Handle(http.MethodPost, "upload", Upload, "0011", nil)
}
