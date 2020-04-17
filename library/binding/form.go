package binding

import (
	"net/http"

	"github.com/pkg/errors"
)

const defaultMemory = 32 * 1024 * 1024

type formBinding struct{}
type formPostBinding struct{}
type formMultipartBinding struct{}

func (f formBinding) Name() string {
	return "form"
}

func (f formBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return errors.WithStack(err)
	}
	if err := mappingByTag(obj, req.Form, "form"); err != nil {
		return err
	}
	return validate(obj)
}

func (f formPostBinding) Name() string {
	return "form-urlencoded"
}

func (f formPostBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return errors.WithStack(err)
	}
	if err := mappingByTag(obj, req.PostForm, "form"); err != nil {
		return err
	}
	return validate(obj)
}

func (f formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (f formMultipartBinding) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return errors.WithStack(err)
	}
	if err := mappingByTag(obj, req.MultipartForm.Value, "form"); err != nil {
		return err
	}
	return validate(obj)
}
