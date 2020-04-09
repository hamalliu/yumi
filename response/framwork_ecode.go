package response

import "yumi/utils/log"

func Success() Status {
	desc := "操作成功"
	return Status{Code: 1, Desc: desc}
}

func InternalError(err error) Status {
	desc := "内部错误"
	log.Error2(err)
	return Status{Code: -500, Desc: desc}
}

//对接通信错误（>1000）（日志级别：Critical，可能遭到通信攻击）
func TokenIsNull() Status {
	desc := "token为空"
	log.Critical2(desc)
	return Status{Code: -1001, Desc: desc}
}

func DecryptError(err error) Status {
	desc := "解密失败"
	log.Critical2(err)
	return Status{Code: -1002, Desc: desc}
}

func SignError(err error) Status {
	desc := "签名错误"
	log.Critical2(err)
	return Status{Code: -1003, Desc: desc}
}

func NoPower(err error) Status {
	desc := "无权访问"
	log.Critical2(err)
	return Status{Code: -1004, Desc: desc}
}

func UntrackedError(err error) Status {
	desc := "内部错误"
	log.Critical2(err)
	return Status{Code: -1005, Desc: desc}
}

func ExpiredSession() Status {
	desc := "会话过期"
	return Status{Code: -1006, Desc: desc}
}
