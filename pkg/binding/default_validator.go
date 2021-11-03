// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ StructValidator = &defaultValidator{}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		//config := &validator.Config{TagName: "binding"}
		v.validate = validator.New()
		f1 := func(fl validator.FieldLevel) bool {
			str := fl.Field().String()
			err := RegexpPhone(str)
			if err != nil {
				return false
			}
			return true
		}
		if err := v.validate.RegisterValidation("phone_number", f1, false); err != nil {
			panic(err)
		}
		f2 := func(fl validator.FieldLevel) bool {
			str := fl.Field().String()
			err := RegexpUser(str)
			if err != nil {
				return false
			}
			return true
		}
		if err := v.validate.RegisterValidation("user_id", f2, false); err != nil {
			panic(err)
		}
		f3 := func(fl validator.FieldLevel) bool {
			str := fl.Field().String()
			err := RegexpPassword(str)
			if err != nil {
				return false
			}
			return true
		}
		if err := v.validate.RegisterValidation("password", f3, false); err != nil {
			panic(err)
		}
		v.validate.SetTagName("binding")
	})
}

// RegexpUser ...
func RegexpUser(str string) error {
	re := regexp.MustCompile(`\w+`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("要求字母，数字或下划线")
	}

	return nil
}

// RegexpPhone ...
func RegexpPhone(str string) error {
	re := regexp.MustCompile(`^1[3-9]\d{9}`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("错误的电话号码")
	}

	return nil
}

// RegexpPassword ...
func RegexpPassword(str string) error {
	if len(str) < 10 || len(str) > 30 {
		return fmt.Errorf("密码字数在10~30之间")
	}

	specailChars := "[~`!@#\\$%\\^&\\*\\(\\)_\\+-=\\{}\\[]\\|\\\\:;\"'<>\\,\\./\\?]"
	alphabet := "[A-Za-z]"
	number := "[0-9]"

	//包含特殊字符和字母
	specialCharAndAlphabet := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", specailChars, alphabet, alphabet, specailChars)
	re := regexp.MustCompile(specialCharAndAlphabet)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	//包含特殊字符和数字
	specialCharAndNumber := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", specailChars, number, number, specailChars)
	re = regexp.MustCompile(specialCharAndNumber)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	//包含数字和字母
	alphabetAndNumber := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", alphabet, number, number, alphabet)
	re = regexp.MustCompile(alphabetAndNumber)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	return fmt.Errorf("密码必须包含特殊字符，数字，字母中的两种以上")
}
