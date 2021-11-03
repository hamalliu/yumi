package toolbox

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	kyExcelCol  = "col"
	kyExcelRow  = "row"
	kyExcelCell = "cell"
)

// ParseExcelToStruct ...
func ParseExcelToStruct(path string, tabIndex int, structItf interface{}, kyExcel string, start int, keys []int) error {
	if kyExcel == "" {
		kyExcel = kyExcelCell
	}
	x, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(structItf)
	tabname := x.GetSheetName(tabIndex)

	switch v.Elem().Kind() {
	case reflect.Struct:
		if err := setSheetToStruct(x, tabname, v.Elem(), v.Elem().Type(), 0, kyExcel); err != nil {
			return err
		}
	case reflect.Slice:
		l := getExcelLen(kyExcel, tabname, x, keys, v)
		l -= start
		mv := reflect.MakeSlice(v.Elem().Type(), l, l)
		for i := 0; i < l; i++ {
			if err := setSheetToStruct(x, tabname, mv.Index(i), mv.Index(i).Type(), i, kyExcel); err != nil {
				return err
			}
		}
		v.Elem().Set(mv)
	default:
		err := errors.New("structItf 参数类型不支持。")
		return err
	}

	return nil
}

func setSheetToStruct(x *excelize.File, tabName string, v reflect.Value, t reflect.Type, index int, kyExcel string) error {
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			if err := setSheetToStruct(x, tabName, v.Field(i), v.Field(i).Type(), index, kyExcel); err != nil {
				return err
			}
		}
		axissstr := t.Field(i).Tag.Get("xls")
		axiss := strings.Split(axissstr, ",")
		axis := ""
		if axissstr == "" {
			continue
		}

		switch kyExcel {
		case kyExcelCol:
			axis = regexp.MustCompile("[a-zA-Z]+").FindString(axiss[0]) + excelRowAdd(regexp.MustCompile("[0-9]+").FindString(axiss[0]), index)
		case kyExcelRow:
			axis = excelColAdd(regexp.MustCompile("[a-zA-Z]+").FindString(axiss[0]), index) + regexp.MustCompile("[0-9]+").FindString(axiss[0])
		case kyExcelCell:
			axis = axiss[index]
		default:
			axis = axiss[index]
		}

		value, _ := x.GetCellValue(tabName, axis)
		if value != "" {
			switch v.Field(i).Kind() {
			case reflect.String:
				if v.Field(i).CanSet() {
					v.Field(i).SetString(value)
				}
			case reflect.Float64:
				f, err := strconv.ParseFloat(value, 64)
				if err != nil {
					errStr := fmt.Sprintf("表格：%s，坐标：%s， 内容：%s只能全是数字!", tabName, axis, value)
					err = errors.New(errStr)
					return err
				}
				if v.Field(i).CanSet() {
					v.Field(i).SetFloat(f)
				}
			case reflect.Int:
				vi, err := strconv.Atoi(value)
				if err != nil {
					fvi, err := strconv.ParseFloat(value, 64)
					if err != nil {
						errStr := fmt.Sprintf("表格：%s，坐标：%s， 内容：%s只能全是整数!", tabName, axis, value)
						err = errors.New(errStr)
						return err
					}
					vi = int(fvi)
				}
				if v.Field(i).CanSet() {
					v.Field(i).SetInt(int64(vi))
				}
			default:
			}
		}
	}

	return nil
}

func getExcelLen(kyExcel, tabname string, x *excelize.File, keys []int, v reflect.Value) (l int) {
	rows, _ := x.GetRows(tabname)

	switch kyExcel {
	case kyExcelCol:
		l = len(rows)
		for i := l - 1; i >= 0; i-- {
			sub := false
			for _, key := range keys {
				if strings.TrimSpace(rows[i][key]) == "" {
					sub = true
					break
				}
			}
			if sub {
				l--
			}
		}
	case kyExcelRow:
		l = len(rows[0])
		for i := l - 1; i >= 0; i-- {
			sub := false
			for _, key := range keys {
				if rows[key][i] == "" {
					sub = true
					break
				}
			}
			if sub {
				l--
			}
		}
	case kyExcelCell:
		l = v.Elem().Len()
	}

	return
}

func excelColAdd(col string, index int) string {
	l := len(col)
	col = strings.ToUpper(col)
	alp := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if strings.Index(alp, string(col[l-1]))+(index%len(alp))+1 == len(alp) {
		col = col[:l-1] + "AA"
	} else {
		col = string(alp[(strings.Index(alp, string(col[l-1]))+index)%(len(alp)-1)])
	}
	return col
}

func excelRowAdd(row string, index int) string {
	rowi, _ := strconv.Atoi(row)
	rowi += index
	return strconv.Itoa(rowi)
}

// OutputStructToExcel 导出到excel
func OutputStructToExcel(x *excelize.File, tabIndex int, structItf interface{}, start int, kyExcel string) error {
	if kyExcel == "" {
		kyExcel = kyExcelCell
	}

	v := reflect.ValueOf(structItf)
	tabname := x.GetSheetName(tabIndex)
	if tabname == "" {
		tabname = fmt.Sprintf("Sheet%d", tabIndex)
		x.NewSheet(tabname)
	}

	switch v.Elem().Kind() {
	case reflect.Struct:
		if err := setStructToSheet(x, tabname, v.Elem(), v.Elem().Type(), 0, kyExcel); err != nil {
			return err
		}
	case reflect.Slice:
		l := v.Elem().Len()
		for i := 0; i < l; i++ {
			ri := i + start
			if err := setStructToSheet(x, tabname, v.Elem().Index(i), v.Elem().Index(i).Type(), ri, kyExcel); err != nil {
				return err
			}
		}
	default:
		err := errors.New("structItf 参数类型不支持。")
		return err
	}

	return nil
}

func setStructToSheet(x *excelize.File, tabName string, v reflect.Value, t reflect.Type, index int, kyExcel string) error {
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			if err := setStructToSheet(x, tabName, v.Field(i), v.Field(i).Type(), index, kyExcel); err != nil {
				return err
			}
		}
		axissstr := t.Field(i).Tag.Get("xls")
		if axissstr == "" {
			continue
		}
		axiss := strings.Split(axissstr, ",")
		axis := ""

		switch kyExcel {
		case kyExcelCol:
			axis = regexp.MustCompile("[a-zA-Z]+").FindString(axiss[0]) + excelRowAdd(regexp.MustCompile("[0-9]+").FindString(axiss[0]), index)
		case kyExcelRow:
			axis = excelColAdd(regexp.MustCompile("[a-zA-Z]+").FindString(axiss[0]), index) + regexp.MustCompile("[0-9]+").FindString(axiss[0])
		case kyExcelCell:
			axis = axiss[index]
		default:
			axis = axiss[index]
		}

		if !v.Field(i).IsZero() {
			if err := x.SetCellValue(tabName, axis, v.Field(i)); err != nil {
				return err
			}
		}
	}

	return nil
}
