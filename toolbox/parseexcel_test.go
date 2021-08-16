package toolbox

import (
	"context"
	"testing"
	"yumi/conf"
	"yumi/pkg/stores/mgoc"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ProdInfo struct {
	Ver   string `json:"ver" bson:"ver"`     // 客户端版本
	DBVer int    `json:"dbver" bson:"dbver"` // 客户端病毒库版本
	TDVer string `json:"tdver" bson:"tdver"` // TD客户端版本
}

type OSInfo struct {
	Name string `json:"name" bson:"name" xls:"H2"` // 操作系统名称
	Ver  string `json:"ver" bson:"ver"`            // 操作系统版本信息
	Arch string `json:"arch" bson:"arch"`          // 32-bit or 64-bit
}

type HWInfo struct {
	Board  string `json:"board" bson:"board"`   // 主板信息
	Memory string `json:"memory" bson:"memory"` // 内存大小
	HDD    string `json:"hdd" bson:"hdd"`       // 硬盘大小
	Video  string `json:"video" bson:"video"`   // 显卡信息
	CPU    string `json:"cpu" bson:"cpu"`
}

type OnlineRequest struct {
}

type User struct {
	Location string `json:"location" bson:"location"`
	Person   string `json:"person" bson:"person" xls:"C2"`
}

type ClientInfo struct {
	Tid      string   `json:"tenant_id" bson:"tenancy_key ,omitempty"` //租户ID
	Hostname string   `json:"hostname" bson:"hostname" xls:"E2"`
	Mac      string   `json:"mac" bson:"mac" xls:"B2"`
	Prod     ProdInfo `json:"prod" bson:"prod"`
	OS       OSInfo   `json:"osinfo" bson:"osinfo" xls:"exist"`
	HW       HWInfo   `json:"hwinfo" bson:"hwinfo"`
	User     User     `json:"user" bson:"user" xls:"exist"`
	OSType   string   `json:"ostype" bson:"ostype" xls:"I2"` //操作系统类型

	ID             string   `json:"_id" bson:"_id" xls:"G2"`                   // 客户端ID
	TenancyKey     string   `json:"tenancy_id" bson:"tenancy_id"`              //租户ID
	ClientIP       string   `json:"client_ip" bson:"client_ip" xls:"A2"`       // 真实IP
	ConnectIP      string   `json:"connect_ip" bson:"connect_ip" xls:"A2"`     // 连接IP connect_ip
	UpgradeVersion string   `json:"upgrade" bson:"upgrade,omitempty" xls:"F2"` // 客户端病毒库版本
	GroupID        string   `json:"group" bson:"group,omitempty" xls:"D2"`     // 客户端所在组ID
	Tags           []string `json:"tags" bson:"tags,omitempty"`

	DefaultTags []string `json:"default_tags" bson:"default_tags,omitempty"` //上线自动贴标签（lili）

	LoginTime int `json:"logintime" bson:"logintime"`

	NextHeartTime int    `json:"next_heart_time" bson:"next_heart_time"` // 时间戳 python中long
	HeartBeatTime int    `json:"heartbeattime" bson:"heartbeattime"`
	Status        int    `json:"status" bson:"status" xls:"K2"`
	PSN           string `json:"psn" bson:"psn"`
	BindPSN       bool   `json:"bind_psn" bson:"bind_psn"`
}

func TestInput(t *testing.T) {
	confMgo := conf.Mongo{
		Dsn: "mongodb://10.34.4.16:27017",
	}

	mgoCli, err := mgoc.New(confMgo.Dsn, confMgo.Options()...)
	if err != nil {
		t.Error(err)
	}

	path := "./topav1.xlsx"
	cs := []*ClientInfo{}
	err = ParseExcelToStruct(path, 1, &cs, kyExcelCol, 1, []int{})
	if err != nil {
		t.Error(err)
	}
	coll := mgoCli.Database("client").Collection("clientinfo")
	docus := []interface{}{}
	for _, v := range cs {

		if v.GroupID == "" {
			v.GroupID = "nogroup"
		}
		docus = append(docus, v)
	}
	ret, err := coll.InsertMany(context.Background(), docus)
	if err != nil {
		t.Error(err)
	}
	t.Log(len(ret.InsertedIDs))
}

func TestDataValidate(t *testing.T) {
	xlsx := excelize.NewFile()
	xlsx.SetCellStr("Sheet1", "A1", "信任方式")
	xlsx.SetCellStr("Sheet1", "B1", "信任特征")
	xlsx.SetCellStr("Sheet1", "C1", "信任程序行为")
	xlsx.SetCellStr("Sheet1", "D1", "是否启用")
	xlsx.SetCellStr("Sheet1", "E1", "备注信息")

	style, err := xlsx.NewStyle(`{"font":{"bold":true,"size":12}, "fill":{"type":"pattern","color":["#FF9933"],"pattern":1}}`)
	if err != nil {
		t.Error(err)
		return
	}
	xlsx.SetCellStyle("Sheet1", "A1", "E1", style)

	dvRange1 := excelize.NewDataValidation(true)
	dvRange1.Sqref = "A2:A9999"
	dvRange1.SetDropList([]string{"path", "sha1"})
	xlsx.AddDataValidation("Sheet1", dvRange1)

	dvRange2 := excelize.NewDataValidation(true)
	dvRange2.Sqref = "C2:C9999"
	dvRange2.SetDropList([]string{"是", "否"})
	xlsx.AddDataValidation("Sheet1", dvRange2)

	dvRange3 := excelize.NewDataValidation(true)
	dvRange3.Sqref = "D2:D9999"
	dvRange3.SetDropList([]string{"启用", "禁用"})
	xlsx.AddDataValidation("Sheet1", dvRange3)

	err = xlsx.SaveAs("c.xlsx")
	if err != nil {
		t.Error(err)
	}
}

type WhiteList struct {
	Data_type  string `bson:"data_type" xls:"A2"`  // 信任方式
	Data_value string `bson:"data_value" xls:"B2"` // 信任特征
	ActionStr  string `bson:"-" xls:"C2"`
	Action     bool   `bson:"action" ` // 信任程序行为
	EnabledStr string `bson:"-" xls:"D2"`
	Enabled    bool   `bson:"enabled"`         // 是否启用
	Remark     string `bson:"remark" xls:"E2"` // 备注信息
}

func TestParseC(t *testing.T) {
	// mulf, err := os.Open("./1.xlsx")
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	ws := []WhiteList{}
	err := ParseExcelToStruct("./1.xlsx", 1, &ws, "col", 1, []int{1, 2, 3, 4})
	if err != nil {
		t.Error(err)
	}
	for i := range ws {
		if ws[i].ActionStr == "是" {
			ws[i].Action = true
		}
		if ws[i].EnabledStr == "启用" {
			ws[i].Enabled = true
		}
	}
	t.Log(ws)
}
