package controller

type Controller struct {
	srve   *Server
	pem    Pemission
	hndlcs HandlerConfs
}

var ctlr Controller

func Init() {
	ctlr.srve = newServer()
	ctlr.pem.UserIdCode = make(map[string]map[string]bool)
	ctlr.hndlcs.handlers = make(map[string]HandlerConf)
}

func GetPemission() *Pemission {
	return &ctlr.pem
}

func GetHandlerConfs() *HandlerConfs {
	return &ctlr.hndlcs
}

func Run() error {
	return ctlr.srve.Run()
}
