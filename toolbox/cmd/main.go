package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"yumi/conf"
	"yumi/pkg/stores/mgoc"
	"yumi/toolbox"
)

type ProdInfo struct {
	Ver   string `json:"ver" bson:"ver" xls:"F2"` // 客户端版本
	DBVer int    `json:"dbver" bson:"dbver"`      // 客户端病毒库版本
	TDVer string `json:"tdver" bson:"tdver"`      // TD客户端版本
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
	TDStatus      int    `josn:"td_status" bson:"td_status"`
	Status        int    `json:"status" bson:"status"`
	StatusDesc    string `json:"-" bson:"-" xls:"K2"`
	PSN           string `json:"psn" bson:"psn"`
	BindPSN       bool   `json:"bind_psn" bson:"bind_psn"`
}

type RETenancy struct {
	KEY           string `json:"key" bson:"_id"`
	UserName      string `json:"username" bson:"username"`
	NickName      string `json:"nickname" bson:"nickname"`
	Address       string `json:"address" bson:"address"`
	Email         string `json:"email" bson:"email"`
	Phone         string `json:"phone" bson:"phone"`
	IsActive      string `json:"is_active" bson:"is_active"`
	RoleName      string `json:"role_name" bson:"role_name"`
	IsSuperman    bool   `json:"is_superman" bson:"is_superman"`
	CreateTime    int64  `json:"create_time" bson:"create_time"`
	PassWord      string `json:"password" bson:"password"`
	CurrTime      int64  `json:"curr_time"`
	Authorization string `json:"-" bson:"authorization"`
}

func LeadIn(dsn, dbName, collectionName, path, tenancyID string) error {
	confMgo := conf.Mongo{
		Dsn: dsn,
	}

	mgoCli, err := mgoc.New(confMgo.Dsn, confMgo.Options()...)
	if err != nil {
		return err
	}

	if tenancyID == "" {
		adminTenancy := RETenancy{}
		err = mgoCli.Database("user_manage").Collection("tenancy_info").FindOne(context.Background(), primitive.M{"default_init_user": true}).Decode(&adminTenancy)
		if err != nil {
			return err
		}
		tenancyID = adminTenancy.KEY
	}

	cs := []ClientInfo{}
	err = toolbox.ParseExcelToStruct(path, 0, &cs, "col", 1, []int{})
	if err != nil {
		return err
	}
	coll := mgoCli.Database(dbName).Collection(collectionName)
	docus := []interface{}{}
	for i := range cs {
		rint := rand.Int()
		cs[i].LoginTime = int(time.Now().AddDate(0, 0, -7).Add(time.Duration(rint%24*5) * time.Hour).Unix())
		cs[i].HeartBeatTime = int(time.Now().AddDate(0, 0, -1).Add(time.Duration(rint%3600) * time.Second).Unix())
		cs[i].TenancyKey = tenancyID
		cs[i].TDStatus = 1

		if cs[i].GroupID == "" {
			cs[i].GroupID = "nogroup"
		}

		if strings.TrimSpace(cs[i].StatusDesc) == "在线" {
			cs[i].Status = 0
		} else if strings.TrimSpace(cs[i].StatusDesc) == "离线" {
			cs[i].Status = 1
		} else {
			return fmt.Errorf("在线状态只能是：在线，离线， 终端ID:%s", cs[i].ID)
		}

		if strings.TrimSpace(cs[i].User.Location) == "" {
			cs[i].User.Location = "t1"
			cs[i].User.Person = "无"
		}

		docus = append(docus, cs[i])
		coll.DeleteOne(context.Background(), primitive.M{"_id": cs[i].ID})
	}

	ret, err := coll.InsertMany(context.Background(), docus)
	if err != nil {
		return err
	}

	fmt.Println("导入条数：", len(ret.InsertedIDs))
	return nil
}

func startFromTerminal() {
REDO:
	fmt.Println("请输入mongo数据库dsn（格式：mongodb://root:Admin_123@10.34.4.16:27017）：")
	dsn := ""
	fmt.Scanln(&dsn)
	fmt.Print("请输入数据库名称和集合名称，使用空格分隔：")
	dbName, collectionName := "", ""
	fmt.Scanln(&dbName, &collectionName)
	fmt.Println("请输入需要导入的excel文件路径：")
GOON:
	path := ""
	fmt.Scanln(&path)

	err := LeadIn(dsn, dbName, collectionName, path, "")
	if err != nil {
		fmt.Println("导入失败，ERROR：", err)
		fmt.Print("输入redo重试，quit退出：")

		do := ""
		fmt.Scanln(&do)
		if do == "redo" {
			goto REDO
		}
	} else {
		fmt.Println("导入成功！")
		fmt.Print("请输入goon继续导入：quit退出：")
		do := ""
		fmt.Scanln(&do)
		if do == "goon" {
			goto GOON
		}
	}
}

type Config struct {
	Dsn            string `json:"dsn"`
	DBName         string `json:"db_name"`
	CollectionName string `json:"collection_name"`
	LeadInPath     string `json:"lead_in_path"`
	TenancyID      string `json:"tenancy_id"`
}

func startFromConfig() {
	path := "./config.json"
	conf := Config{}
	f, err := os.Open(path)
	if err != nil {
		goto ERROR
	}
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		goto ERROR
	}
	err = LeadIn(conf.Dsn, conf.DBName, conf.CollectionName, conf.LeadInPath, conf.TenancyID)
	if err != nil {
		goto ERROR
	}
	fmt.Println("导入成功！")
	fmt.Println("按Enter键退出！")
	fmt.Scanln()
	return
ERROR:
	fmt.Println("导入失败！ERROR：", err)
	fmt.Println("按Enter键退出！")
	fmt.Scanln()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover  error: ", err)
		}
	}()

	a := 0
	print(5 / a)
}
