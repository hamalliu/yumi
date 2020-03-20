package response

import (
	"errors"
	"fmt"

	"yumi/utils/internal_error"
	"yumi/utils/log"
)

type Status struct {
	code   int    //状态码
	desc   string //状态描述
	Detail string //状态详情
}

const (
	SolutionGeneral = "请联系管理员或稍后再试"
)

//操作成功
func Success() Status {
	return Status{code: 1000, desc: "操作成功"}
}

func Error(detail error) Status {
	var targetErr *internal_error.Error
	if errors.As(detail, &targetErr) {
		//后端程序错误
		return internalError()
	} else {
		//用户输入错误
		return illegalParameter(detail)
	}
}

//对接通信错误（>1000）（日志级别：Critical，可能遭到通信攻击）
func TokenIsNull(detail error) Status {
	desc := "token为空"
	log.Critical2(fmt.Sprintf("%s:%v", desc, detail))
	return Status{code: 1001, desc: desc}
}
func DecryptError(detail error) Status {
	desc := "解密失败"
	log.Critical2(fmt.Sprintf("%s:%v", desc, detail))
	return Status{code: 1002, desc: desc}
}
func SignError(detail error) Status {
	desc := "签名错误"
	log.Critical2(fmt.Sprintf("%s:%v", desc, detail))
	return Status{code: 1003, desc: desc}
}
func NoPower(detail error) Status {
	desc := "无权访问"
	log.Critical2(fmt.Sprintf("%s:%v", desc, detail))
	return Status{code: 1004, desc: desc}
}
func ExpiredSession(detail error) Status {
	desc := "会话过期"
	//log.Critical2(fmt.Sprintf("%s:%v", desc, detail))
	return Status{code: 1005, desc: desc}
}

//后端程序错误（2000）（日志级别：Error，管理员应及时查看该日志，改进系统性能）
//内部错误不将详情返给前端
func internalError() Status {
	desc := "内部错误"
	return Status{code: 2000, desc: desc, Detail: SolutionGeneral}
}

//用户输入错误（3000）（日志级别：Info，管理员应及时查看该日志，改进系统性能）
func illegalParameter(detail error) Status {
	desc := "参数不合法"
	log.Info2(fmt.Sprintf("%s:%v", desc, detail))
	return Status{code: 3000, desc: desc, Detail: fmt.Sprintf("%v", detail)}
}
