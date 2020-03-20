package pretreatment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"yumi/response"
)

func UrlUnmarshal(reqUrl *url.URL, reqJ interface{}) response.Status {
	var (
		reqJv = reflect.ValueOf(reqJ)
	)
	if reqJv.Kind() != reflect.Ptr || reqJv.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("reqJ必须为结构体指针"))
	}

	fl := reqJv.NumField()

	//验证参数格式
	for i := 0; i < fl; i++ {
		val := reqUrl.Query().Get(strings.ToLower(reqJv.Type().Field(i).Name))
		if val == "" {
			continue
		}
		switch reqJv.Kind() {
		case reflect.Float32, reflect.Float64:
			if fval, err := strconv.ParseFloat(val, 64); err != nil {
				err := fmt.Errorf("参数：%s 格式错误", strings.ToLower(reqJv.Type().Field(i).Name))
				return response.Error(err)
			} else {
				reqJv.Field(i).SetFloat(fval)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uval, err := strconv.ParseUint(val, 0, 64); err != nil {
				err := fmt.Errorf("参数：%s 格式错误", strings.ToLower(reqJv.Type().Field(i).Name))
				return response.Error(err)
			} else {
				reqJv.Field(i).SetUint(uval)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if ival, err := strconv.ParseInt(val, 0, 64); err != nil {
				err := fmt.Errorf("参数：%s 格式错误", strings.ToLower(reqJv.Type().Field(i).Name))
				return response.Error(err)
			} else {
				reqJv.Field(i).SetInt(ival)
			}
		case reflect.Bool:
			if bval, err := strconv.ParseBool(val); err != nil {
				err := fmt.Errorf("参数：%s 格式错误", strings.ToLower(reqJv.Type().Field(i).Name))
				return response.Error(err)
			} else {
				reqJv.Field(i).SetBool(bval)
			}
		case reflect.String:
			reqJv.Field(i).SetString(val)
		}
	}

	//检查
	for i := 0; i < fl; i++ {
		if err := feildCheck(reqJv.Field(i), reqJv.Type().Field(i)); err != nil {
			return response.Error(err)
		}
	}

	return response.Success()
}

func BodyUnmarshal(data io.ReadCloser, reqJ interface{}) response.Status {
	reqJv := reflect.ValueOf(reqJ)
	if reqJv.Kind() != reflect.Ptr || reqJv.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("reqJ必须为结构体指针"))
	}

	if err := json.NewDecoder(data).Decode(reqJ); err != nil {
		return response.Error(err)
	}

	//检查
	reqJv = reflect.ValueOf(reqJ)
	if err := checkStruct(reqJv); err != nil {
		return response.Error(err)
	}

	return response.Success()
}

func checkStruct(value reflect.Value) error {
	fl := value.Elem().NumField()
	for i := 0; i < fl; i++ {
		if value.Elem().Field(i).Kind() == reflect.Struct {
			return checkStruct(value.Field(i))
		}

		if value.Elem().Field(i).Kind() == reflect.Slice {
			sl := value.Elem().Field(i).Len()
			for si := 0; si < sl; si++ {
				if value.Elem().Field(i).Index(si).Kind() == reflect.Ptr &&
					value.Elem().Field(i).Index(si).Elem().Kind() == reflect.Struct {
					return checkStruct(value.Elem().Field(i).Index(si))
				} else if value.Elem().Field(i).Index(si).Elem().Kind() == reflect.Struct {
					return checkStruct(value.Elem().Field(i).Index(si).Addr())
				}
			}
		}

		if err := feildCheck(value.Elem().Field(i), value.Elem().Type().Field(i)); err != nil {
			return err
		}
	}

	return nil
}

func feildCheck(feild reflect.Value, feildType reflect.StructField) error {
	tag := feildType.Tag
	tagMap := getTag(tag)
	feildDesc := tagMap["desc"]
	if tagMap["notzero"] == "notzero" {
		if feild.IsZero() {
			return fmt.Errorf("%s， 不能为零值", feildDesc)
		}
	}
	feildStr := fmt.Sprint(feild.Interface())

	if err := strfmtCheck(tagMap["strfmt"], feildStr); err != nil {
		return fmt.Errorf("%s， %s", feildDesc, err.Error())
	}

	if tagMap["size"] != "" {
		size, _ := strconv.Atoi(tagMap["size"])
		if len([]rune(feildStr)) > size {
			return fmt.Errorf("%s，不能超过%d个字符", feildDesc, size)
		}
	}

	if tagMap["regexp"] != "" {
		re := regexp.MustCompile(tagMap["regexp"])
		if string(re.Find([]byte(feildStr))) != feildStr {
			return fmt.Errorf("%s， 格式错误", feildDesc)
		}
	}

	return nil
}

func getTag(tag reflect.StructTag) map[string]string {
	tmaps := make(map[string]string)
	tagStr := tag.Get("check")
	tagStr = strings.ReplaceAll(tagStr, " ", "")
	tags := strings.Split(tagStr, ";")
	for i := range tags {
		vs := strings.Split(tags[i], ":")
		if len(vs) >= 2 {
			tmaps[vs[0]] = vs[1]
		} else {
			tmaps[vs[0]] = vs[0]
		}
	}

	return tmaps
}

func strfmtCheck(strfmt, str string) error {
	if strfmt == "" {
		return nil
	}
	switch strfmt {
	case "number":
		return CheckNumber(str)
	case "user":
		return CheckUser(str)
	case "phone":
		return CheckPhone(str)
	case "ip":
		return CheckIp(str)
	default:
		return fmt.Errorf(`不支持的strfmt`)
	}
}
