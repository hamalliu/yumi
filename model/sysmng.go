package model

import (
	"time"
)

//账号
type Account struct {
	Id           int64     `db:"id" json:"id"`
	User         string    `db:"user" json:"user"`                 //用户名(用于登录)
	Name         string    `db:"name" json:"name"`                 //真实姓名
	Code         string    `db:"code" json:"code"`                 //账号编码
	Mobile       string    `db:"mobile" json:"mobile"`             //电话号码
	Password     string    `db:"password" json:"password"`         //密码
	RegisterTime time.Time `db:"registertime" json:"registertime"` //注册时间
	Status       string    `db:"status" json:"status"`             //状态：关闭，启用

	Operator    string    `db:"operator" json:"operator"`       //修改人
	OperateTime time.Time `db:"operatetime" json:"operatetime"` //修改时间

	Checked bool `db:"-" json:"_checked"` //备注
}

//角色
type Role struct {
	Id     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`     //名称
	Code   string `db:"code" json:"code"`     //编码
	Status string `db:"status" json:"status"` //状态：关闭，启用

	Operator    string    `db:"operator" json:"operator"`       //修改人
	OperateTime time.Time `db:"operatetime" json:"operatetime"` //更新时间

	Checked bool `db:"-" json:"_checked"` //是否被选中
}

//用户角色关联
type AcctRoles struct {
	Id       int64  `db:"id" json:"id"`
	AcctCode string `db:"usercode" json:"usercode"` //账号编码
	RoleCode string `db:"rolecode" json:"rolecode"` //角色编码
}

//菜单
type Menu struct {
	Id           int64  `db:"id" json:"id"`
	ParentName   string `db:"parentname" json:"parentname"`     //父级名称
	ParentCode   string `db:"parentcode" json:"parentcode"`     //父级编码
	Name         string `db:"name" json:"name"`                 //名称
	Code         string `db:"code" json:"code"`                 //编码
	Route        string `db:"route" json:"route"`               //英文名
	Params       string `db:"params" json:"params"`             //参数
	Type         string `db:"type" json:"type"`                 //类型（菜单，功能）
	DisplayOrder int64  `db:"displayorder" json:"displayorder"` //显示顺序
	Status       string `db:"status" json:"status"`             //状态：关闭，启用

	CurSubCode  uint `db:"cursubcode" json:"-"`  //最新子菜单编码（以便推算下一个子菜单编码）
	CurFuncCode uint `db:"curfunccode" json:"-"` //最新菜单功能编码（以便推算下一个子菜单编码）

	Operator    string    `db:"operator" json:"operator"`       //操作者
	OperateTime time.Time `db:"operatetime" json:"operatetime"` //操作时间
}

//账号菜单关联
type AcctMenus struct {
	Id       int64
	AcctCode string `db:"acctcode" json:"acctcode"`
	MenuCode string `db:"menucode" json:"menucode"`
}

//角色菜单关联
type RoleMenus struct {
	Id       int64
	RoleCode string `db:"rolecode" json:"rolecode"`
	MenuCode string `db:"menucode" json:"menucode"`
}

//选中菜单树
type SelectMenuTree struct {
	Id       int64            `json:"id"`
	Name     string           `json:"name"`     //名称
	Code     string           `json:"code"`     //编码
	Select   bool             `json:"select"`   //是否选中
	Expand   bool             `json:"expand"`   //是否展开
	Children []SelectMenuTree `json:"children"` //子级
}

//勾选菜单树
type CheckedMenuTree struct {
	Id       int64             `json:"id"`
	Name     string            `json:"name"`     //名称
	Code     string            `json:"code"`     //编码
	Checked  bool              `json:"checked"`  //是否选中
	Expand   bool              `json:"expand"`   //是否展开
	Children []CheckedMenuTree `json:"children"` //子级
}

//权限树
type Power struct {
	ParentName string  `json:"parentname"` //父级名称
	ParentCode string  `json:"parentcode"` //父级编码
	Name       string  `json:"name"`       //名称
	Code       string  `json:"code"`       //编码
	Route      string  `json:"route"`      //英文名
	Params     string  `json:"params"`     //参数
	Children   []Power `json:"children"`   //子菜单
	Functions  []Power `json:"functions"`  //功能
}

//系统数据字典
type SysDataDict struct {
	Id          int64     `db:"id" json:"id"`
	Name        string    `db:"dictname" json:"dictname"`       //字典类型名称
	Code        string    `db:"dictcode" json:"dictcode"`       //字典类型编码
	Remark      string    `db:"remark" json:"remark"`           //备注
	OperateTime time.Time `db:"operatetime" json:"operatetime"` //更新时间
	Operator    string    `db:"operator" json:"operator"`       //操作人
	RowNumber   int       `db:"rownumber" json:"rownumber"`     //行号
	Checked     bool      `db:"-" json:"_checked"`              //是否选中
}

//系统数据字典item
type SysDataDictItem struct {
	Id          int64     `db:"id" json:"id"`
	Pptr        string    `db:"pptr" json:"pptr"`
	Name        string    `db:"name" json:"name"`
	Code        string    `db:"code" json:"code"`
	OperateTime time.Time `db:"operatetime" json:"operatetime"` //更新时间
	Operator    string    `db:"operator" json:"operator"`       //操作人
	RowNumber   int       `db:"rownumber" json:"rownumber"`     //行号
	Checked     bool      `db:"-" json:"_checked"`              //是否选中
}

//更新记录
type UpdateRecordsJson struct {
	Id          int    `db:"id" json:"id"`
	Table       string `db:"table" json:"table"`
	Request     string `db:"request" json:"request"`
	RequestBody string `db:"requestbody" json:"requestbody"`
	BeforeData  string `db:"beforedata" json:"beforedata"`
	UserId      string `db:"userid" json:"userid"`
	Operator    string `db:"operator" json:"operator"`
	OperateTime string `db:"operatetime" json:"operatetime"`
}

//删除记录
type DeleteRecordsJson struct {
	Id          int    `db:"id" json:"id"`
	Table       string `db:"table" json:"table"`
	Request     string `db:"request" json:"request"`
	RequestBody string `db:"requestbody" json:"requestbody"`
	BeforeData  string `db:"beforedata" json:"beforedata"`
	UserId      string `db:"userid" json:"userid"`
	Operator    string `db:"operator" json:"operator"`
	OperateTime string `db:"operatetime" json:"operatetime"`
}
