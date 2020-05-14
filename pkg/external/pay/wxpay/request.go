package wxpay

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"yumi/pkg/external/pay/internal"
)

var reqRcrd *internal.ReqRcrd

func init() {
	reqRcrd = internal.NewReqRcrd("wxpay_req_rcrds")
}

func request(respBody interface{}, method, url string, reqBody interface{}, tr *http.Transport) ([]byte, error) {
	if reflect.ValueOf(respBody).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("参数respBody必须为指针")
	}
	if reflect.ValueOf(reqBody).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("参数reqBody必须为指针")
	}

	body, err := xml.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ctxtimeout, _ := context.WithTimeout(ctx, time.Second*15)
	reqtx, err := http.NewRequestWithContext(ctxtimeout, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	//发起请求
	cli := http.Client{}
	if tr != nil {
		cli.Transport = tr
	}
	resp, err := cli.Do(reqtx)
	if err != nil {
		if err == http.ErrHandlerTimeout {
			return nil, fmt.Errorf("请求超时")
		}
		return nil, err
	}

	bs, _ := ioutil.ReadAll(resp.Body)
	//解析返回数据
	if err := xml.Unmarshal(bs, respBody); err != nil {
		return bs, fmt.Errorf("%s", string(bs))
	}

	//记录请求
	reqRcrd.AddRcrds(method, url, reqtx.Header, body, resp.Header, bs)

	return bs, nil
}
