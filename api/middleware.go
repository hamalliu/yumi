package api

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"yumi/conf"
	"yumi/controller"
	"yumi/response"
	"yumi/session"
	"yumi/utils"
	"yumi/utils/log"
)

//解密
func Decrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if !controller.GetHandlerConfs().Get(mux.CurrentRoute(req).GetName()).GetReqEncrypt() {
			next.ServeHTTP(resp, req)
			return
		}

		var (
			timeStamp = req.Header.Get("timestamp")
			user      = req.Header.Get("user")
			s         session.Session

			ok  bool
			err error
		)

		if s, ok = session.GetUser(user); !ok {
			response.Response(resp, req, response.ExpiredSession(nil), nil)
		}

		cryted, err := ioutil.ReadAll(req.Body)
		if err != nil {
			response.Response(resp, req, response.DecryptError(err), nil)
		}
		key := utils.MD5([]byte(s.Token + timeStamp))
		body, err := utils.AesDecrypt(string(cryted), key)
		if err != nil {
			response.Response(resp, req, response.DecryptError(err), nil)
		}
		if err = req.Write(bytes.NewBufferString(body)); err != nil {
			response.Response(resp, req, response.DecryptError(err), nil)
		}

		next.ServeHTTP(resp, req)
	})
}

//
func DefenseRepaly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		//TODO
		return
	})
}

//验签

//验权
func PemissionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		paternCode := mux.CurrentRoute(req).GetName()
		if controller.GetHandlerConfs().Get(paternCode).GetSkipPemissionAuth() {
			next.ServeHTTP(resp, req)
			return
		}
		if !controller.GetPemission().HavePower(req.Header.Get("user"), paternCode) {
			response.Response(resp, req, response.NoPower(nil), nil)
			return
		}
		next.ServeHTTP(resp, req)
		return
	})
}

func DebugLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if conf.Get().Environment != conf.EnvDebug {
			next.ServeHTTP(resp, req)
			return
		}

		log.Debug("req:", req.URL.String())
		log.Debug("body:", req.Body)
	})
}
