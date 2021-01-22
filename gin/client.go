package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"yumi/pkg/render"
)

const timeformat = "2006-01-02 15:04:05.999999"

// Client is
type Client struct {
	Addr string
	User string
	Sid  string
}

// New a client.
func New(path, user, pwd string) Client {

	return Client{}
}

//GetJSON 以json格式交互的get请求
func (cli Client) GetJSON(path string, urlParams map[string]string) render.JSON {
	urlVals := url.Values{}
	for k, v := range urlParams {
		urlVals.Set(k, v)
	}
	requrl := fmt.Sprintf("http://%s%s", cli.Addr, path)
	if urlVals.Encode() != "" {
		requrl = fmt.Sprintf("%s?%s", requrl, urlVals.Encode())
	}

	req, err := http.NewRequest(http.MethodGet, requrl, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("req_%s: %s", time.Now().Format(timeformat), requrl)
	var httpCli http.Client
	resp, err := httpCli.Do(req)
	if err != nil {
		panic(err)
	}
	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp_%s: %s", time.Now().Format(timeformat), string(respBodyBytes))

	respJ := render.JSON{}
	err = json.Unmarshal(respBodyBytes, &respJ)
	if err != nil {
		panic(err)
	}

	return respJ
}

//PostJSON 以json格式交互的post请求
func (cli Client) PostJSON(path string, body interface{}, urlParams map[string]string) render.JSON {
	urlVals := url.Values{}

	for k, v := range urlParams {
		urlVals.Set(k, v)
	}
	requrl := fmt.Sprintf("http://%s%s", cli.Addr, path)
	if urlVals.Encode() != "" {
		requrl = fmt.Sprintf("%s?%s", requrl, urlVals.Encode())
	}

	var (
		bodyByte []byte
		err      error
	)
	if body != nil {
		bodyByte, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	req, err := http.NewRequest(http.MethodPost, requrl, bytes.NewBuffer(bodyByte))
	if err != nil {
		panic(err)
	}
	fmt.Printf("req_%s: %s", time.Now().Format(timeformat), requrl)
	var httpCli http.Client
	resp, err := httpCli.Do(req)
	if err != nil {
		panic(err)
	}
	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp_%s: %s", time.Now().Format(timeformat), string(respBodyBytes))

	respJ := render.JSON{}
	err = json.Unmarshal(respBodyBytes, &respJ)
	if err != nil {
		panic(err)
	}

	return respJ
}
