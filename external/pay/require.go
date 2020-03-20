package pay

import (
	"fmt"
	"reflect"
)

func CheckRequire(obj interface{}) error {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("obj类型错误")
	}

	fldl := t.NumField()
	for i := 0; i < fldl; i++ {
		if t.Field(i).Tag.Get("require") == "true" {
			if v.Field(i).IsZero() {
				return fmt.Errorf("字段：%s不能为零值", t.Field(i).Name)
			}
		}
	}

	return nil
}
