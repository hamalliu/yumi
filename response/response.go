package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"yumi/conf"
	"yumi/consts"
	"yumi/controller"
	"yumi/utils"
	"yumi/utils/log"
)

type Data struct {
	Code     int         `json:"code"`
	Desc     string      `json:"desc"`
	Business interface{} `json:"business"`
}

func (m Data) Bytes() []byte {
	respBytes, _ := json.Marshal(m)
	return respBytes
}

func (m Data) String() string {
	respBytes, _ := json.Marshal(m)
	return string(respBytes)
}

type PageItem struct {
	Page  int         `json:"page"`
	Line  int         `json:"line"`
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

func Response(resp http.ResponseWriter, req *http.Request, err error, item interface{}) {
	data := &Data{}
	status := Status{}

	if err == nil {
		status = Success()
	} else {
		if errors.As(err, &Status{}) {
			status = err.(Status)
		} else {
			status = UntrackedError(err)
		}
	}
	data.Code = status.Code
	data.Desc = status.Desc
	data.Business = item

	timeStamp := time.Now().UnixNano()
	resp.Header().Set(consts.HeaderTimestamp, fmt.Sprint(timeStamp))

	if controller.GetHandlerConfs().Get(mux.CurrentRoute(req).GetName()).GetRespEncrypt() {
		resp.Header().Set(consts.HeaderEncrypt, "true")

		token := req.Header.Get("token")
		if token == "" {
			status = TokenIsNull()
			data.Code = status.Code
			data.Desc = status.Desc
			goto ERROR

		} else {
			key := utils.MD5([]byte(fmt.Sprint(timeStamp) + token))
			if encode, err := utils.AesEncrypt(data.String(), key); err != nil {
				status = InternalError(err)
				data.Code = status.Code
				data.Desc = status.Desc
				goto ERROR

			} else {
				_, _ = resp.Write([]byte(encode))
			}
		}

	ERROR:
		_, _ = resp.Write(data.Bytes())
	} else {
		resp.Header().Set(consts.HeaderEncrypt, "false")

		_, _ = resp.Write(data.Bytes())
	}

	if conf.EnvDebug == conf.Get().Environment {
		log.Debug("resp:", data.String())
	}
}
