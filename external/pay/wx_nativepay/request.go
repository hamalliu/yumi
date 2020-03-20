package wx_nativepay

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

func request(respBody interface{}, method, url string, reqBody interface{}) ([]byte, error) {
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
	var cli http.Client
	resp, err := cli.Do(reqtx)
	if err != nil {
		if err == http.ErrHandlerTimeout {
			return nil, fmt.Errorf("请求超时")
		}
		return nil, err
	}

	//解析返回数据
	if err := xml.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return ioutil.ReadAll(resp.Body)
	}

	return nil, nil
}
