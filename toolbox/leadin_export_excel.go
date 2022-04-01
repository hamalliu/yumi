package toolbox

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	ContentTypeCol  = "col"
	ContentTypeRow  = "row"
	ContentTypeCell = "cell"
)

// ExcelModel excel模板
type ExcelModel struct {
	Header         map[string]string
	Start          int
	ContentType    string
	ContentRequire map[string]ContentRequire
}

// ContentRequire
type ContentRequire struct {
	DropList   map[string]interface{}
	RangeSqref string
}

// Marshal 打包当前模板
func (m ExcelModel) Marshal(x *excelize.File) error {
	tabName := "Sheet1"
	for k, v := range m.Header {
		if err := x.SetCellStr(tabName, k, v); err != nil {
			return err
		}
	}
	for _, v := range m.ContentRequire {
		if len(v.DropList) > 0 && v.RangeSqref != "" {
			dvRange := excelize.NewDataValidation(true)
			dvRange.Sqref = v.RangeSqref
			var dl []string
			for k := range v.DropList {
				dl = append(dl, k)
			}
			dvRange.SetDropList(dl)
			dvRange.SetError(excelize.DataValidationErrorStyleStop, "", "")
			x.AddDataValidation(tabName, dvRange)
		}
	}

	return nil
}

// IsCurrentModel 导入的是否是当前模板
func (m ExcelModel) IsCurrentModel(x *excelize.File, tabName string) bool {
	for k, v := range m.Header {
		if val, err := x.GetCellValue(tabName, k); err != nil || val != v {
			return false
		}
	}

	return true
}

// ParseExcelPathToStruct 解析excel到对象
func (m ExcelModel) ParseExcelToObject(x *excelize.File, obj interface{}, tabIndex int) error {
	v := reflect.ValueOf(obj)
	tabname := x.GetSheetName(tabIndex)

	if !m.IsCurrentModel(x, tabname) {
		return errors.New("不是当前导入模板")
	}

	switch v.Elem().Kind() {
	case reflect.Struct:
		if err := m.setSheetToStruct(x, tabname, v.Elem(), v.Elem().Type(), 0); err != nil {
			return err
		}
	case reflect.Slice:
		l := getExcelLen(m.ContentType, tabname, x, v)
		l -= m.Start
		mv := reflect.MakeSlice(v.Elem().Type(), l, l)
		for i := 0; i < l; i++ {
			if err := m.setSheetToStruct(x, tabname, mv.Index(i), mv.Index(i).Type(), i); err != nil {
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

// ParseExcelPathToObject 解析excel文件的路径到对象
func (m ExcelModel) ParseExcelPathToObject(path string, obj interface{}, tabIndex int) error {
	x, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	return m.ParseExcelToObject(x, obj, tabIndex)
}

// ParseExcelReaderToObject 解析excel的reader到对象
func (m ExcelModel) ParseExcelReaderToObject(read io.Reader, obj interface{}, tabIndex int) error {
	x, err := excelize.OpenReader(read)
	if err != nil {
		return err
	}

	return m.ParseExcelToObject(x, obj, tabIndex)
}

func (m ExcelModel) setSheetToStruct(x *excelize.File, tabName string, v reflect.Value, t reflect.Type, index int) error {
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			if err := m.setSheetToStruct(x, tabName, v.Field(i), v.Field(i).Type(), index); err != nil {
				return err
			}
		}
		axissstr := t.Field(i).Tag.Get("xls")
		axiss := strings.Split(axissstr, ",")
		axis := ""
		if axissstr == "" || axissstr == "-" {
			continue
		}

		switch m.ContentType {
		case ContentTypeCol:
			axis = regexp.MustCompile("[a-zA-Z]+").FindString(axiss[0]) + excelRowAdd(regexp.MustCompile("[0-9]+").FindString(axiss[0]), index)
		case ContentTypeRow:
			axis = excelColAdd(regexp.MustCompile("[a-zA-Z]+").FindString(axiss[0]), index) + regexp.MustCompile("[0-9]+").FindString(axiss[0])
		case ContentTypeCell:
			axis = axiss[index]
		default:
			axis = axiss[index]
		}

		value, _ := x.GetCellValue(tabName, axis)
		if value != "" {
			if dlv, ok := m.ContentRequire[axissstr]; ok {
				convertVal, exist := dlv.DropList[value]
				if !exist {
					return fmt.Errorf("表格：%s，坐标：%s， 内容：%s不是数据验证中的数据!", tabName, axis, value)
				}
				v.Field(i).Set(reflect.ValueOf(convertVal))
			} else {
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
	}

	return nil
}

// ExportToExcelObject 导出对象
type ExportToExcelObject interface {
	GetExportValueByHeaderName(hn string) (value interface{})
}

// ParseObjectToExcelNoReflect 不使用反射解析导出对象到excel
func (m ExcelModel) ParseObjectToExcelWithoutReflect(x *excelize.File, tabName string, objs []ExportToExcelObject) (err error) {
	index := m.Start
	for _, v := range objs {
		index++
		for k := range m.Header {
			axis := k
			switch m.ContentType {
			case ContentTypeCol:
				axis = regexp.MustCompile("[a-zA-Z]+").FindString(axis) + strconv.Itoa(index)
			case ContentTypeRow:
				axis = excelColAdd(regexp.MustCompile("[a-zA-Z]+").FindString(axis), index) +
					regexp.MustCompile("[0-9]+").FindString(axis)
			default:
				return errors.New("no support content_type")
			}

			value := v.GetExportValueByHeaderName(k)
			if err = x.SetCellValue(tabName, axis, value); err != nil {
				return fmt.Errorf("parse object to excel without reflect:%w", err)
			}
		}
	}

	return nil
}

// ParseObjectToExcel 解析对象到excel
func (m ExcelModel) ParseObjectToExcel(x *excelize.File, obj interface{}, tabIndex int) error {
	v := reflect.ValueOf(obj)
	tabname := x.GetSheetName(tabIndex)
	if tabname == "" {
		tabname = fmt.Sprintf("Sheet%d", tabIndex)
		x.NewSheet(tabname)
	}

	switch v.Elem().Kind() {
	case reflect.Struct:
		if err := m.setStructToSheet(x, tabname, v.Elem(), v.Elem().Type(), 0); err != nil {
			return err
		}
	case reflect.Slice:
		l := v.Elem().Len()
		for i := 0; i < l; i++ {
			ri := i + m.Start
			if err := m.setStructToSheet(x, tabname, v.Elem().Index(i), v.Elem().Index(i).Type(), ri); err != nil {
				return err
			}
		}
	default:
		err := errors.New("structItf 参数类型不支持。")
		return err
	}

	return nil
}

func (m ExcelModel) setStructToSheet(x *excelize.File, tabName string, v reflect.Value, t reflect.Type, index int) error {
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			if err := m.setStructToSheet(x, tabName, v.Field(i), v.Field(i).Type(), index); err != nil {
				return err
			}
		}
		axissstr := t.Field(i).Tag.Get("xls")
		if axissstr == "" || axissstr == "-" {
			continue
		}
		axiss := strings.Split(axissstr, ",")
		axis := ""

		switch m.ContentType {
		case ContentTypeCol:
			axis = regexp.MustCompile("[a-zA-Z]+").FindString(axiss[0]) + excelRowAdd(regexp.MustCompile("[0-9]+").FindString(axiss[0]), index)
		case ContentTypeRow:
			axis = excelColAdd(regexp.MustCompile("[a-zA-Z]+").FindString(axiss[0]), index) + regexp.MustCompile("[0-9]+").FindString(axiss[0])
		case ContentTypeCell:
			axis = axiss[index]
		default:
			axis = axiss[index]
		}

		if !v.Field(i).IsZero() {
			x.SetCellValue(tabName, axis, v.Field(i))
		}
	}

	return nil
}

func getExcelLen(kyExcel, tabname string, x *excelize.File, v reflect.Value) (l int) {
	rows, _ := x.GetRows(tabname)

	switch kyExcel {
	case ContentTypeCol:
		l = len(rows)
	case ContentTypeRow:
		l = len(rows[0])
	case ContentTypeCell:
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
