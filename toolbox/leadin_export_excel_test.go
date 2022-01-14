package toolbox

import (
	"testing"

	"github.com/xuri/excelize/v2"
)


func TestParseExcelToObject(t *testing.T) {
	x, err := excelize.OpenFile("./1.xlsx")
	if err != nil {
		t.Error(err)
	}

	t.Log(x.GetCellFormula("Sheet1", "D2"))
}
