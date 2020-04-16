package controller

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime"
	"time"

	"github.com/gorilla/mux"

	"yumi/conf"
	"yumi/consts"
	"yumi/utils/log"
)

type Server struct {
	re *mux.Router
}

func newServer() (srvr *Server) {
	srvr = new(Server)
	srvr.re = mux.NewRouter()
	srvr.re.Use(CORSMiddleware, Recovery)

	return
}

func (m *Server) Group(pattern string, handlers []interface{}) {
	var muxHandlers []mux.MiddlewareFunc
	for i := range handlers {
		muxHandlers = append(muxHandlers, handlers[i].(func(http.Handler) http.Handler))
	}

	m.re.PathPrefix(pattern).Subrouter().Use(muxHandlers...)
}

func (m *Server) Handle(httpMethod, pattern string, handler interface{}, patternCode string) {
	m.re.HandleFunc(pattern, handler.(func(http.ResponseWriter, *http.Request))).Methods(httpMethod).Name(patternCode)
}

func (m *Server) Run() error {
	httpsrv := &http.Server{
		Handler:      m.re,
		Addr:         conf.Get().Addr,
		ReadTimeout:  time.Duration(conf.Get().ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.Get().WriteTimeout) * time.Second,
	}

	return httpsrv.ListenAndServe()
}

//Cors跨域中间件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		header := "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
		header = fmt.Sprintf("%s, %s", header, consts.GetHeaders())
		w.Header().Set("Access-Control-Allow-Headers", header)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("options return OK"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			var rawReq []byte
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				if request != nil {
					rawReq, _ = httputil.DumpRequest(request, false)
				}
				pl := fmt.Sprintf("http call panic: %s\n%v\n%s\n", string(rawReq), err, buf)
				log.Error(pl)
				writer.WriteHeader(500)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
