package toolbox

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

// 异常进程黑名单
var _AbnormalProcessBlackListLeadinModel = ExcelModel{
	Header: map[string]string{
		"A1": "名称",
		"B1": "进程路径",
		"C1": "状态",
	},
	ContentRequire: map[string]ContentRequire{
		"C1": {
			DropList: map[string]interface{}{
				"启用": 1,
				"禁用": 2,
			},
			RangeSqref: "C2:C9999",
		},
	},
	ContentType: ContentTypeCol,
	Start:       1,
}

type AbnormalProcessBlackList struct {
	Name        string `json:"name" bson:"name" xls:"A2"`                 // 名称
	ProcessPath string `json:"process_path" bson:"process_path" xls:"B2"` // 进程路径
	Status      int8   `json:"status" bson:"status" xls:"C2"`             // 1-启用 2-禁用
}

func (m *AbnormalProcessBlackList) GetExportValueByHeaderName(hn string) (value interface{}) {
	switch hn {
	case "A1":
		return m.Name
	case "B1":
		return m.ProcessPath
	case "C1":
		switch m.Status {
		case 1:
			return "启用"
		case 2:
			return "禁用"
		default:
			return ""
		}
	default:
		return nil
	}
}

func TestParseExcelToObject(t *testing.T) {
	x, err := excelize.OpenFile("./1.xlsx")
	if err != nil {
		t.Error(err)
	}

	t.Log(x.GetCellFormula("Sheet1", "D2"))
}

func TestParseObjectToExcelWithoutReflect(t *testing.T) {
	excelize.NewFile()
	x := excelize.NewFile()
	objs := []ExportToExcelObject{}
	objs = append(objs, &AbnormalProcessBlackList{
		Name:        "1",
		ProcessPath: "C:/aaa",
		Status:      1,
	})
	if err := _AbnormalProcessBlackListLeadinModel.Marshal(x); err != nil {
		t.Error(err)
	}
	if err := _AbnormalProcessBlackListLeadinModel.ParseObjectToExcelWithoutReflect(x, "Sheet2", objs); err != nil {
		t.Error(err)
		return
	}
	if err := x.SaveAs("2.xlsx"); err != nil {
		t.Error(err)
		return
	}
}
