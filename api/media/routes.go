package media

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	mr := r.Group("media")

	mr.Handle(http.MethodPost, "", "/uploadmultipart", UploadMultipart)
	mr.Handle(http.MethodPost, "", "/upload", Upload)
}
