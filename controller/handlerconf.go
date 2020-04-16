package controller

type HandlerConf struct {
	//请求加密
	reqEncrypt bool
	//返回加密
	respEncrypt bool
	//跳过权限认证
	skipPemissionAuth bool
}

func NewHandlerConf(reqEncrypt, respEncrypt, skipPemissionAuth bool) *HandlerConf {
	return &HandlerConf{
		reqEncrypt:        reqEncrypt,
		respEncrypt:       respEncrypt,
		skipPemissionAuth: skipPemissionAuth,
	}
}

func (hndl HandlerConf) GetReqEncrypt() bool {
	return hndl.reqEncrypt
}

func (hndl HandlerConf) GetRespEncrypt() bool {
	return hndl.respEncrypt
}

func (hndl HandlerConf) GetSkipPemissionAuth() bool {
	return hndl.skipPemissionAuth
}

type HandlerConfs struct {
	handlers map[string] /*patternCode*/ HandlerConf
}

func (hndl *HandlerConfs) add(pattern string, hdlCnf *HandlerConf) {
	//设置默认值
	if hdlCnf == nil {
		hdlCnf = &HandlerConf{
			reqEncrypt:  true,
			respEncrypt: true,
		}
	}

	hndl.handlers[pattern] = *hdlCnf
}

func (hndl *HandlerConfs) Get(pattern string) HandlerConf {
	return hndl.handlers[pattern]
}
