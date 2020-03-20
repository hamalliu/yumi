package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"yumi/conf"
	"yumi/controller"
	"yumi/utils"
	"yumi/utils/log"
)

type Data struct {
	Code     int         `json:"code"`
	Desc     string      `json:"desc"`
	Detail   string      `json:"detail"`
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

func Response(resp http.ResponseWriter, req *http.Request, status Status, item interface{}) {
	timeStamp := time.Now().UnixNano()
	resp.Header().Set("timestamp", fmt.Sprint(timeStamp))
	data := &Data{status.code, status.desc, status.Detail, item}

	if controller.GetHandlerConfs().Get(mux.CurrentRoute(req).GetName()).GetRespEncrypt() {
		resp.Header().Set("encrypt", "true")

		token := req.Header.Get("token")
		if token == "" {
			err := TokenIsNull(nil)
			_, _ = resp.Write(Data{Code: err.code, Desc: err.desc, Detail: err.Detail}.Bytes())
			return
		}
		key := utils.MD5([]byte(fmt.Sprint(timeStamp) + token))
		encode, err := utils.AesEncrypt(data.String(), key)
		if err != nil {
			errStatus := Error(err)
			_, _ = resp.Write(Data{Code: errStatus.code, Desc: errStatus.desc, Detail: errStatus.Detail}.Bytes())
		} else {
			_, _ = resp.Write([]byte(encode))
		}
	} else {
		resp.Header().Set("encrypt", "false")

		_, _ = resp.Write(data.Bytes())
	}

	if conf.EnvDebug == conf.Get().Environment {
		log.Debug("resp:", data.String())
	}
}
