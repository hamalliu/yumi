package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"yumi/conf"
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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("options return OK"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

//恢复panic中间件
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				stack := stack(3)
				httpRequest, _ := httputil.DumpRequest(request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				log.Error(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s", timeFormat(time.Now()), err, stack))

				if !brokenPipe {
					writer.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned apiclient
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}
