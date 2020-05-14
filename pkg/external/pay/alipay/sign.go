package alipay

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

//验证回调签名
func NotifyVerify(order interface{}, signData string, pubkey string) error {
	orderStr := ""

	switch v := reflect.ValueOf(order); v.Kind() {
	case reflect.String:
		orderStr = v.String()

	case reflect.Map:
		orderStr = notifyMap2Sign(v.Interface())

	case reflect.Struct:
		orderStr = notifyStruct2Sign(v.Interface())

	case reflect.Ptr:
		orderStr = notifyStruct2Sign(v.Elem().Interface())

	default:
		panic("params type not supported")
	}

	return rsa2Verify(orderStr, signData, pubkey)
}

func RespVerify(order interface{}, signData string, pubkey string) error {
	orderStr := ""

	switch v := reflect.ValueOf(order); v.Kind() {
	case reflect.String:
		orderStr = v.String()

	case reflect.Map:
		orderStr = respMap2Sign(v.Interface())

	case reflect.Struct:
		orderStr = respStruct2Sign(v.Interface())

	case reflect.Ptr:
		orderStr = respStruct2Sign(v.Elem().Interface())

	default:
		panic("params type not supported")
	}

	return rsa2Verify(orderStr, signData, pubkey)
}

func BuildSign(order interface{}, key string) (string, error) {
	orderStr := ""

	switch v := reflect.ValueOf(order); v.Kind() {
	case reflect.String:
		orderStr = v.String()

	case reflect.Map:
		orderStr = map2Sign(v.Interface())

	case reflect.Struct:
		orderStr = struct2Sign(v.Interface())

	case reflect.Ptr:
		orderStr = struct2Sign(v.Elem().Interface())

	default:
		panic("params type not supported")
	}

	if b, err := rsa2Sign([]byte(orderStr), key); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(b), nil
	}
}

func respStruct2Sign(content interface{}) string {
	tempMap := make(map[string]interface{})
	v := reflect.ValueOf(content)
	t := reflect.TypeOf(content)
	l := t.NumField()
	for i := 0; i < l; i++ {
		if v.Field(i).IsZero() {
			continue
		}
		k := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		vv := fmt.Sprintf("%v", v.Field(i).Interface())
		k = strings.ReplaceAll(k, " ", "")
		vv = strings.ReplaceAll(vv, " ", "")

		if k != "-" && k != "" && k != "sign" {
			tempMap[k] = vv
		}
	}

	tempString, _ := json.Marshal(&tempMap)
	return string(tempString)
}

func respMap2Sign(content interface{}) string {
	tempMap := make(map[string]interface{})
	switch v := content.(type) {
	case map[string]interface{}:
		for k := range v {
			k = strings.ReplaceAll(k, " ", "")
			if k == "sign" {
				continue
			}
			tempMap[k] = v[k]
		}
	}
	tempString, _ := json.Marshal(&tempMap)
	return string(tempString)
}

func notifyStruct2Sign(content interface{}) string {
	var tempArr []string
	temString := ""
	v := reflect.ValueOf(content)
	t := reflect.TypeOf(content)
	l := t.NumField()
	for i := 0; i < l; i++ {
		if v.Field(i).IsZero() {
			continue
		}
		k := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		vv := fmt.Sprintf("%v", v.Field(i).Interface())
		k = strings.ReplaceAll(k, " ", "")
		vv = strings.ReplaceAll(vv, " ", "")

		if k != "-" && k != "" && k != "sign" && k != "sign_type" {
			tempArr = append(tempArr, k+"="+vv)
		}
	}

	sort.Strings(tempArr)
	first := true
	for _, v := range tempArr {
		if first {
			temString = v
			first = false
		} else {
			temString = fmt.Sprintf("%s&%s", temString, v)
		}
	}

	return temString
}

func notifyMap2Sign(content interface{}) string {
	tempStr := ""
	switch v := content.(type) {
	case map[string]interface{}:
		var buf bytes.Buffer
		keys := make([]string, 0, len(v))

		for k := range v {
			k = strings.ReplaceAll(k, " ", "")
			if k == "sign" || k == "sign_type" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if v[k] == "" {
				continue
			}
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}

			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(fmt.Sprintf("%v", v[k]))
		}
		tempStr = buf.String()
	}
	return tempStr
}

func map2Sign(content interface{}) string {
	tempStr := ""
	switch v := content.(type) {
	case map[string]interface{}:
		var buf bytes.Buffer
		keys := make([]string, 0, len(v))

		for k := range v {
			k = strings.ReplaceAll(k, " ", "")
			if k == "sign" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if v[k] == "" {
				continue
			}
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}

			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(fmt.Sprintf("%v", v[k]))
		}
		tempStr = buf.String()
	}
	return tempStr
}

func struct2Sign(content interface{}) string {
	var tempArr []string
	temString := ""
	v := reflect.ValueOf(content)
	t := reflect.TypeOf(content)
	l := t.NumField()
	for i := 0; i < l; i++ {
		if v.Field(i).IsZero() {
			continue
		}
		k := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		vv := fmt.Sprintf("%v", v.Field(i).Interface())
		k = strings.ReplaceAll(k, " ", "")
		vv = strings.ReplaceAll(vv, " ", "")

		if k != "-" && k != "" && k != "sign" {
			tempArr = append(tempArr, k+"="+vv)
		}
	}

	sort.Strings(tempArr)
	first := true
	for _, v := range tempArr {
		if first {
			temString = v
			first = false
		} else {
			temString = fmt.Sprintf("%s&%s", temString, v)
		}
	}

	return temString
}

func rsa2Verify(data, signStr, pubkey string) error {
	sign, err := base64.StdEncoding.DecodeString(signStr)
	if err != nil {
		return err
	}

	block, _ := pem.Decode([]byte(pubkey))
	if block == nil {
		return errors.New("block为空")
	}
	pubix, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pub := pubix.(*rsa.PublicKey)

	hashed := sha256.Sum256([]byte(data))

	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], sign)
}

func rsa2Sign(data []byte, prikey string) ([]byte, error) {
	block2, _ := pem.Decode([]byte(prikey))
	if block2 == nil {
		return nil, errors.New("block为空")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block2.Bytes)
	if err != nil {
		return nil, err
	}
	p := priv.(*rsa.PrivateKey)

	rng := rand.Reader
	hashed := sha256.Sum256(data)

	return rsa.SignPKCS1v15(rng, p, crypto.SHA256, hashed[:])
}
