package gin

import (
	"net/http"
)

//Run 运行服务器
func Run(server *http.Server) error {
	_ = server.ListenAndServe()
	return nil
}

//RunTLS 运行https服务器
func RunTLS(server *http.Server) error {
	//TODO
	return nil
}

//RunUnix 运行unix域服务器
func RunUnix(server *http.Server) error {
	//TODO
	return nil
}
