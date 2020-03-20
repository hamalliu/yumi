package controller

import (
	"context"
	"net/http"
)

func ContextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

func ContextSet(r *http.Request, key, val interface{}) *http.Request {
	if val == nil {
		return r
	}

	return r.WithContext(context.WithValue(r.Context(), key, val))
}
