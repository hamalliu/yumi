package ali_pagepay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	body, err := json.Marshal(reqBody)
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
	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return ioutil.ReadAll(resp.Body)
	}

	return nil, nil
}

func ParseQuery(rawQuery string) (ReqNotify, error) {
	notify := ReqNotify{}

	vals, err := url.ParseQuery(rawQuery)
	if err != nil {
		return notify, err
	}

	t := reflect.TypeOf(notify)
	v := reflect.ValueOf(&notify)
	fl := t.NumField()

	for i := 0; i < fl; i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		if jsonTag != "-" &&
			len(vals[jsonTag]) != 0 {
			v.Elem().Field(i).SetString(vals[jsonTag][0])
		}
	}

	return notify, nil
}
