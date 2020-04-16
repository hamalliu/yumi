package ecodes

import "yumi/utils/log"

/**
 *框架级别错误码（-500-500）
 */

func Success() Status {
	desc := "操作成功"
	return Status{Code: 101, Desc: desc}
}

func InternalError(err error) Status {
	desc := "服务器错误"
	log.Error2(err)
	return Status{Code: -101, Desc: desc}
}

func TokenIsNull() Status {
	desc := "token为空"
	log.Critical2(desc)
	return Status{Code: -102, Desc: desc}
}

func DecryptError(err error) Status {
	desc := "解密失败"
	log.Critical2(err)
	return Status{Code: -102, Desc: desc}
}

func SignError(err error) Status {
	desc := "签名错误"
	log.Critical2(err)
	return Status{Code: -103, Desc: desc}
}

func NoPower(err error) Status {
	desc := "无权访问"
	log.Critical2(err)
	return Status{Code: -104, Desc: desc}
}

func UntrackedError(err error) Status {
	desc := "未追踪错误"
	log.Critical2(err)
	return Status{Code: -105, Desc: desc}
}

func ExpiredSession() Status {
	desc := "会话过期"
	return Status{Code: -106, Desc: desc}
}

//参数预处理错误
func PretreatmentError(err error) Status {
	return Status{Code: -107, Desc: err.Error()}
}

func FormDataTooLarge(err error) Status {
	desc := "form数据太大"
	log.Error2(err)
	return Status{Code: -108, Desc: desc}
}
