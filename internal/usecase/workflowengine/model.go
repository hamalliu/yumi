package workflowengine

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
)

const (
	cronLogPath = "cronWf.log"
)

const (
	RvwTypeRvw       = "rvw"
	RvwtypeOrd       = "ord"
	RvwtypeGroup     = "group"
	RvwtypeAssiOrd   = "assiord"
	RvwtypeAssiGroup = "assigroup"
)

//流程和节点的状态
const (
	StatusNodeInstanceUndone     = "undone" //未完成
	StatusNodeInstanceStatusDone = "done"   //已完成
)

const (
	WfinsStatusTypeSet        = "set"        //流程实例置为完成状态
	WfinsStatusTypePendingSet = "pendingset" //（多分支情况）结束本分支，如果其他分支已结束，流程实例置为完成状态
)

const (
	WorkflowStatusMaking  = "定义中"
	WorkflowStatusRelease = "已发布"
)

type Config struct {
	ResetDb  bool       `json:"resetdb"`
	DbConfig gdb.Config `json:"dbconfig"`
}

type Model struct {
	db    *gorm.DB
	wfMux sync.Mutex

	//工作间(工作间决定流程流出方向)
	workspace map[string]GetOutflow

	//钩子
	hook map[string]Hook

	//拦截
	intercept map[string]Intercept

	//流程节点map/流程map
	nds map[string]map[string]*Node
	wfs map[string]*Workflow
	//流程实例状态
	wfinsStat map[string] /*value*/ string /*type*/

	actors map[string] /*roleid*/ []actor
}

type GetOutflow func(wsId []uint) (string, error)

type Intercept func(instance WorkflowInstance) (bool, error)

type Hook func(instance WorkflowInstance) error

func New(conf Config) (*Model, error) {
	var (
		m Model
		f *os.File

		err error
	)

	if m.db, err = gdb.New(conf.DbConfig); err != nil {
		return nil, fmt.Errorf("数据库连接失败，%s", err.Error())
	}

	m.hook = make(map[string]Hook)
	m.intercept = make(map[string]Intercept)
	m.workspace = make(map[string]GetOutflow)

	if err = m.load(); err != nil {
		return nil, fmt.Errorf("加载流程失败，%s", err.Error())
	}

	if f, err = os.OpenFile(cronLogPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644); err != nil {
		return nil, err
	}
	loger := cron.PrintfLogger(log.New(f, "cron: ", log.LstdFlags))
	crn := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(loger)))
	if _, err = crn.AddFunc("0 0 0 * * ?", m.canceledWfins); err != nil {
		return nil, err
	}
	crn.Start()

	return nil, nil
}

func (m *Model) load() error {
	m.nds = make(map[string]map[string]*Node)
	m.wfs = make(map[string]*Workflow)
	m.wfinsStat = make(map[string]string)

	var (
		wfs      []Workflow
		wfinsSts []WorkflowInstanceStatus
		rls      []role

		err error
	)

	//加载流程
	if err = m.db.Where("status = ?", WorkflowStatusRelease).Find(&wfs).Error; err != nil {
		return fmt.Errorf("查询流程定义错误，%s", err.Error())
	}
	for i := range wfs {
		if err = m.db.Model(&wfs[i]).Related(&wfs[i].Nodes).Error; err != nil {
			return fmt.Errorf("查询流程节点错误，%s", err.Error())
		}
	}
	for wi := range wfs {
		m.wfs[wfs[wi].Code] = &wfs[wi]

		m.nds[wfs[wi].Code] = make(map[string]*Node)
		for ni := range wfs[wi].Nodes {
			wfs[wi].Nodes[ni].RoleMap = parseNodeRoles(wfs[wi].Nodes[ni].Roles)
			wfs[wi].Nodes[ni].OutflowMap = parseOutflow(wfs[wi].Nodes[ni].Outflow)

			m.nds[wfs[wi].Code][wfs[wi].Nodes[ni].Code] = &wfs[wi].Nodes[ni]
		}
	}

	//加载流程实例状态
	if err = m.db.Find(&wfinsSts).Error; err != nil {
		return err
	}
	for i := range wfinsSts {
		m.wfinsStat[wfinsSts[i].Vaule] = wfinsSts[i].Type
	}

	//加载执行者
	if err = m.db.Find(&rls).Error; err != nil {
		return err
	}
	for i := range rls {
		if err = m.db.Related(&rls[i].Actors).Error; err != nil {
			return err
		}
		m.actors[fmt.Sprintf("%d", rls[i].ID)] = rls[i].Actors
	}

	return nil
}

func (m *Model) RegisteHook(hook interface{}) error {
	hookv := reflect.ValueOf(hook)
	if hookv.Kind() != reflect.Ptr {
		return fmt.Errorf("hook必须为指针")
	}

	for _, wv := range m.nds {
		for _, nv := range wv {
			if nv.Hook != "" {
				if hookv.MethodByName(nv.Hook).Kind() == reflect.Invalid ||
					hookv.MethodByName(nv.Hook).Kind().String() != "func(instance WorkflowInstance) error" {
					return fmt.Errorf("流程%s,节点%s,钩子函数%s不存在", nv.WfCode, nv.Code, nv.Hook)
				} else {
					m.hook[nv.Hook] = hookv.MethodByName(nv.Hook).Interface().(func(instance WorkflowInstance) error)
				}
			}
		}
	}

	return nil
}

func (m *Model) RegisteInterceptor(interceptor interface{}) error {
	interceptorv := reflect.ValueOf(interceptor)

	if interceptorv.Kind() != reflect.Ptr {
		return fmt.Errorf("interceptor必须为指针")
	}

	for _, wv := range m.nds {
		for _, nv := range wv {
			if nv.Hook != "" {
				if interceptorv.MethodByName(nv.Intercept).Kind() == reflect.Invalid ||
					interceptorv.MethodByName(nv.Intercept).Kind().String() != "func(instance WorkflowInstance) (bool, error)" {
					return fmt.Errorf("流程%s,节点%s,拦截函数%s不存在", nv.WfCode, nv.Code, nv.Intercept)
				} else {
					m.intercept[nv.Intercept] = interceptorv.MethodByName(nv.Intercept).Interface().(func(instance WorkflowInstance) (bool, error))
				}
			}
		}
	}

	return nil
}

func (m *Model) RegisteWorkSpace(ws interface{}) error {
	wsv := reflect.ValueOf(ws)

	if wsv.Kind() != reflect.Ptr {
		return fmt.Errorf("interceptor必须为指针")
	}

	for _, wv := range m.nds {
		for _, nv := range wv {
			if wsv.MethodByName(nv.Intercept).Kind() == reflect.Invalid ||
				wsv.MethodByName(nv.Intercept).Kind().String() != "func(wsId []uint) (string, error)" {
				return fmt.Errorf("流程%s,节点%s,工作间%s不存在", nv.WfCode, nv.Code, nv.Intercept)
			} else {
				m.workspace[nv.WsName] = wsv.MethodByName(nv.WsName).Interface().(func(wsId []uint) (string, error))
			}
		}
	}

	return nil
}

func parseNodeRoles(roleStr string) (roleMap map[string][]uint) {
	var (
		rprs = strings.Split(roleStr, ";")
	)

	for i := range rprs {
		roleIdsStr := strings.Split(strings.Split(rprs[i], ":")[1], ",")
		for ri := range roleIdsStr {
			roleId, _ := strconv.Atoi(roleIdsStr[ri])
			roleMap[strings.Split(rprs[i], ":")[0]] = append(roleMap[strings.Split(rprs[i], ":")[0]], uint(roleId))
		}

	}

	return
}

func parseOutflow(outflow string) (outflowMap map[string][]string) {
	var (
		outflows = strings.Split(outflow, ";")
	)

	for i := range outflows {
		nodeCodesStr := strings.Split(strings.Split(outflows[i], ":")[1], ",")
		for ri := range nodeCodesStr {
			outflowMap[strings.Split(outflows[i], ":")[0]] = append(outflowMap[strings.Split(outflows[i], ":")[0]], nodeCodesStr[ri])
		}
	}

	return
}
