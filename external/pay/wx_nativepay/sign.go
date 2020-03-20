package wx_nativepay

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

//构建签名
func Buildsign(order interface{}, bizkey string) string {
	orderStr := ""

	switch v := reflect.ValueOf(order); v.Kind() {
	case reflect.String:
		orderStr = v.String()

	case reflect.Map:
		orderStr = map2Sgin(v.Interface(), bizkey)

	case reflect.Struct:
		orderStr = struct2Sign(v.Interface(), bizkey)

	case reflect.Ptr:
		orderStr = struct2Sign(v.Elem().Interface(), bizkey)

	default:
		panic("params type not supported")
	}

	//生成md5签名
	md5ctx := md5.New()
	md5ctx.Write([]byte(orderStr))
	return strings.ToUpper(hex.EncodeToString(md5ctx.Sum(nil)))
}

func map2Sgin(content interface{}, bizKey string) (str string) {
	switch v := content.(type) {
	case map[string]string:
		var buf bytes.Buffer
		keys := make([]string, 0, len(v))

		for k := range v {
			k = strings.ReplaceAll(k, " ", "")
			if k == "sign" ||
				k == "" {
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

			buf.WriteString(fmt.Sprintf("%s=%s", k, v[k]))
		}
		buf.WriteString(fmt.Sprintf("&key=%s", bizKey))
		str = buf.String()
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

			buf.WriteString(fmt.Sprintf("%s=%v", k, v[k]))
		}
		buf.WriteString(fmt.Sprintf("&key=%s", bizKey))
		str = buf.String()
	}
	return str
}

func struct2Sign(content interface{}, bizKey string) string {
	var tempArr []string
	temString := ""
	v := reflect.ValueOf(content)
	t := reflect.TypeOf(content)
	l := t.NumField()
	for i := 0; i < l; i++ {
		if v.Field(i).IsZero() {
			continue
		}
		k := strings.Split(t.Field(i).Tag.Get("xml"), ",")[0]
		vv := fmt.Sprintf("%v", v.Field(i).Interface())
		k = strings.ReplaceAll(k, " ", "")
		vv = strings.ReplaceAll(vv, " ", "")

		if k != "-" && k != "" && k != "sign" {
			tempArr = append(tempArr, k+"="+vv+"&")
		}
	}

	sort.Strings(tempArr)
	for _, v := range tempArr {
		temString = temString + v
	}

	temString = temString + "key=" + bizKey
	return temString
}
