package workflowengine

import (
	"sync"
	"time"
)

//执行者
type Actor struct {
	Deleted bool
	ID      uint
	Name    string //执行者姓名
}

//流程定义
type Workflow struct {
	nodeMux sync.Mutex `json:"-"`

	ID           uint   `json:"id"`
	Code         string `json:"code"`         //流程编码
	Name         string `json:"name"`         //流程名称
	FirstNode    string `json:"firstnode"`    //发起node
	ParentWfCode string `json:"parentwfcode"` //父级流程code
	Status       string `json:"status"`       //定义中，已发布

	Nodes []Node `json:"nodes"`
}

//流程定义节点
type Node struct {
	ID             uint   `json:"id"`
	Code           string `json:"code"`           //节点编码
	WfCode         string `json:"wfcode"`         //流程编码
	Name           string `json:"name"`           //节点名称
	Pahse          string `json:"pahse"`          //阶段名称
	WsName         string `json:"wsname"`         //工作间名称
	Outflow        string `json:"outflow"`        //格式（outflow:nodecode,nodecode;outflow:nodecode,nodecode）
	Intercepter    string `json:"intercepter"`    //拦截函数
	Hook           string `json:"hooks"`          //钩子函数
	RvwType        string `json:"rvwtype"`        //审批类型(order,group)
	ReferenceNodes string `json:"referencenodes"` //参考nodeinstance（以获取该节点执行者）
	Roles          string `json:"roles"`          //执行者（com>dep>role:1,2,3）（节点类型为rvw的时候不能为空）

	OutflowMap map[string] /*outflow*/ []string/*nodecode*/ `json:"-"`
	RoleMap    map[string] /*rolepath*/ []uint/*roleid*/ `json:"-"`

	WorkflowID uint `json:"workflowid"` //流程id
}

//流程实例
type WorkflowInstance struct {
	ID            uint      `json:"id"`
	Year          string    `json:"year"`          //分表标准字段（暂定以年为标准）
	Name          string    `json:"name"`          //流程名称
	Abstract      string    `json:"abstract"`      //流程摘要
	Code          string    `json:"code"`          //流程编码
	Status        string    `json:"status"`        //流程状态（与业务强相关）
	StatusReason  string    `json:"statusreason"`  //流程状态原因（与业务强相关）
	StartDate     time.Time `json:"startdate"`     //申请时间
	EndDate       time.Time `json:"enddate"`       //结束时间
	UpdateDate    time.Time `json:"updatedate"`    //更新时间
	ParentWfinsID uint      `json:"parentwfinsid"` //父级流程实例id

	RolePath            string `json:"rolepath"`            //发起人角色路径（com>dep>role）
	ActorID             uint   `json:"actorid"`             //发起人id
	CurRolePath         string `json:"currolepath"`         //当前执行者角色路径（com>dep>role）
	CurActorID          uint   `json:"curactorid"`          //当前执行者id
	FirstNodeINstanceID uint   `json:"firstnodeinstanceid"` //第一节点实例id

	SubWfinss        []WorkflowInstance                      `json:"subwfinss"`
	NodeInstances    []NodeInstance                          `json:"nodeinstances"`
	NodeInstancesMap map[string] /*code*/ []NodeInstance     `json:"-"`
	SubWfinssMap     map[string] /*code*/ []WorkflowInstance `json:"-"`
}

//流程实例节点
type NodeInstance struct {
	ID                 uint      `json:"id"`
	NdinsSerialNum     uint      `json:"ndinsserialnum"` //同一个节点的节点实例序号
	NdinsActors        string    ``                      //同一个节点的所有执行者
	Pahse              string    `json:"pahse"`          //阶段名称
	WfCode             string    `json:"wfcode"`         //流程编码
	NodeCode           string    `json:"nodecode"`       //节点编码
	NodeName           string    `json:"nodename"`       //节点名称
	StartDate          time.Time `json:"startdate"`      //申请时间
	EndDate            time.Time `json:"enddate"`        //结束时间
	Deadline           time.Time `json:"deadline"`       //截止日期（到了截止日期自动启动推动函数）
	RvwType            string    `json:"rvwtype"`        //审批类型(rvw,ordrvw,grouprvw,assiordrvw,assigrouprvw)
	Status             string    `json:"status"`         //状态
	WsName             string    `json:"wsname"`         //工作间名称
	WsId               uint      `json:"wsid"`           //工作间id
	RolePath           string    `json:"rolepath"`       //角色路径（com>dep>role）
	ActorID            uint      `json:"actorid"`        //执行者id
	UpdateDate         time.Time `json:"updatedate"`     //更新时间
	ParNodeInstanceID  uint      //父节点实例id
	WorkflowInstanceID uint      `json:"workflowinstanceid"` //流程实例id

	Context    string `json:"context"`    //上下文（供不同类型的节点，存储自己格式的内容）
	AssiActors string `json:"assiactors"` //指派执行者
}

//流程实例状态
type WorkflowInstanceStatus struct {
	ID    uint
	Type  string ``
	Vaule string ``
}

//移交任务记录
type moveTaskRecord struct {
	ID          uint
	WfInsID     uint
	AccountID   uint
	BeAccountID uint
}
