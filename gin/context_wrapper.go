package gin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"yumi/gin/header"
	"yumi/gin/valuer"
	"yumi/pkg/binding"
	"yumi/pkg/codec"
	"yumi/pkg/log"
	"yumi/pkg/status"
)

func (c *Context) marshalJSON(data interface{}, err error) (int, []byte, error) {
	c.Error = err

	if data == nil {
		data = struct{}{}
	}

	s := status.OK()
	if err != nil {
		originErr := errors.Unwrap(err)
		ss, ok := originErr.(*status.Status)
		if ok {
			ss.WithError(err)
			s = ss
		} else {
			s = status.Unknown().WithError(err)
		}
		log.Error(s.Error())
	}

	respObj := struct {
		Code    int         `json:"code"`
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{}
	respObj.Code = int(s.Code())
	respObj.Status = s.Code().String()
	respObj.Message = s.Message(c.Request.Header.Get("Accept-Language"))
	respObj.Data = data

	bs, err := json.Marshal(respObj)
	return s.Code().HTTPStatus(), bs, err
}

// WriteJSON json序列化通用返回结构，并写入连接中
func (c *Context) WriteJSON(data interface{}, err error) {
	code, bs, err := c.marshalJSON(data, err)
	if err != nil {
		panic(err)
	}
	c.Bytes(code, "application/json; charset=utf-8", bs)
}

// WriteEncryptJSON json序列化通用返回结构，并加密再写入连接中
func (c *Context) WriteEncryptJSON(data interface{}, err error) {
	code, bs, err := c.marshalJSON(data, err)
	if err != nil {
		panic(err)
	}
	bs, err = c.EncryptBytes(bs)
	if err != nil {
		panic(err)
	}

	header.SetBodyEncrypt(c.Writer)
	c.Bytes(code, "application/json; charset=utf-8", bs)
}

// DecryptBodyAndBind 解密请求中的body再通过默认绑定器，绑定的指定结构体
func (c *Context) DecryptBodyAndBind(obj interface{}) error {
	err := c.DecryptBody()
	if err != nil {
		return err
	}

	return c.Bind(obj)
}

// DecryptBodyAndBindWith 解密请求中的body再通过指定的绑定器，绑定到指定结构体
func (c *Context) DecryptBodyAndBindWith(obj interface{}, b binding.Binding) error {
	err := c.DecryptBody()
	if err != nil {
		return err
	}

	return c.BindWith(obj, b)
}

// DecryptBody 解密请求中的body
func (c *Context) DecryptBody() error {
	if header.BodyEncrypt(c.Request) != "true" {
		return errors.New("no encryption in the request")
	}

	if c.Request.Body == nil {
		return nil
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}
	decryptBodyBytes, err := codec.EcbDecrypt(c.Get(valuer.KeySecret).Bytes(), bodyBytes)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(decryptBodyBytes)
	c.Request.Body = io.NopCloser(buf)
	return nil
}

// EncryptBytes 加密指定的bytes
func (c *Context) EncryptBytes(bs []byte) ([]byte, error) {
	encryptBodyBytes, err := codec.EcbEncrypt(c.Get(valuer.KeySecret).Bytes(), bs)
	if err != nil {
		return nil, err
	}
	return encryptBodyBytes, nil
}
