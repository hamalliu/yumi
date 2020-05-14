package wxpay

import (
	"fmt"
	"reflect"
	"strings"
)

func BuildPrameter(xmlStruct interface{}) string {
	var tempArr []string
	temString := ""
	v := reflect.ValueOf(xmlStruct)
	t := reflect.TypeOf(xmlStruct)
	l := t.NumField()
	for i := 0; i < l; i++ {
		if v.Field(i).IsZero() {
			continue
		}
		k := strings.Split(t.Field(i).Tag.Get("xml"), ",")[0]
		vv := fmt.Sprintf("%v", v.Field(i).Interface())
		k = strings.ReplaceAll(k, " ", "")
		vv = strings.ReplaceAll(vv, " ", "")

		if k != "-" && k != "" {
			tempArr = append(tempArr, k+"="+vv)
		}
	}

	first := true
	for _, v := range tempArr {
		if first {
			temString = temString + v
			first = false
		} else {
			temString = temString + "&" + v
		}
	}

	return temString
}
