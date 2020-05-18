package memory

import (
	"fmt"
	"strconv"
	"sync"

	"yumi/internal/db"
	"yumi/model"
	"yumi/pkg/conf"
)

const (
	MenuItem     = "菜单项"
	FunctionItem = "功能项"
)

type DataMenu struct {
	mux  sync.Mutex
	root *MenuTree
}

type MenuTree struct {
	Id           int64  `json:"id"`
	ParentName   string `json:"parentname"`   //父级名称
	ParentCode   string `json:"parentcode"`   //父级编码
	Name         string `json:"name"`         //名称
	Code         string `json:"code"`         //编码
	Route        string `json:"route"`        //英文名
	Params       string `json:"params"`       //参数
	Type         string `json:"type"`         //类型
	DisplayOrder int64  `json:"displayorder"` //显示顺序
	Status       string `json:"status"`       //状态
	CurSubCode   uint   `json:"-"`            //最新子菜单编码
	CurFuncCode  uint   `json:"-"`            //最新菜单功能编码

	Children    []*MenuTree          `json:"children"`
	childrenMap map[string]*MenuTree `json:"-"`
}

//构建菜单树
func NewDataMenu() *DataMenu {
	mt, parentMenu := newMenuTree()
	buildMenuTree(parentMenu, mt.root)
	mt.root.CurSubCode = uint(len(mt.root.Children))

	return mt
}

func newMenuTree() (*DataMenu, map[string][]*MenuTree) {
	var err error
	mt := new(DataMenu)
	mt.root = new(MenuTree)

	mt.root.Name = conf.Get().SysName
	mt.root.Code = "0"
	mt.root.childrenMap = make(map[string]*MenuTree)

	menus := []model.Menu{}
	if menus, err = db.Menu().GetAllMenus(); err != nil {
		panic(err)
	}
	parentMenu := make(map[string][]*MenuTree)
	for i := range menus {
		menu := new(MenuTree)
		menu.Id = menus[i].Id
		menu.ParentName = menus[i].ParentName
		menu.ParentCode = menus[i].ParentCode
		menu.Name = menus[i].Name
		menu.Code = menus[i].Code
		menu.Route = menus[i].Route
		menu.Params = menus[i].Params
		menu.Type = menus[i].Type
		menu.DisplayOrder = menus[i].DisplayOrder
		menu.Status = menus[i].Status
		menu.CurSubCode = menus[i].CurSubCode
		menu.CurFuncCode = menus[i].CurFuncCode

		mt.root.childrenMap[menu.Code] = menu
		parentMenu[menu.ParentCode] = append(parentMenu[menu.ParentCode], menu)
	}

	return mt, parentMenu
}

func buildMenuTree(parentMenu map[string][]*MenuTree, menu *MenuTree) {
	for i := range parentMenu[menu.Code] {
		buildMenuTree(parentMenu, parentMenu[menu.Code][i])
		menu.Children = append(menu.Children, parentMenu[menu.Code][i])
	}
}

func reload() {
	m = NewDataMenu()
}

var m *DataMenu

func Init() {
	reload()
}

func Menu() *DataMenu {
	return m
}

func (mt *MenuTree) correctCurCode() {
	maxCode := ""
	if mt.Type == MenuItem {
		for i := range mt.Children {
			if maxCode < mt.Children[i].Code[len(mt.Children[i].Code)-3:] {
				maxCode = mt.Children[i].Code[len(mt.Children[i].Code)-3:]
			}
		}
		curCode, _ := strconv.ParseUint(maxCode, 10, 32)
		mt.CurSubCode = uint(curCode)
	} else {
		for i := range mt.Children {
			if maxCode < mt.Children[i].Code[len(mt.Children[i].Code)-2:] {
				maxCode = mt.Children[i].Code[len(mt.Children[i].Code)-2:]
			}
		}
		curCode, _ := strconv.ParseUint(maxCode, 10, 32)
		mt.CurFuncCode = uint(curCode)
	}
}

func (m *DataMenu) Add(menu model.Menu) error {
	var (
		err error
	)
	m.mux.Lock()
	defer m.mux.Unlock()

	if menu.ParentCode == "0" {
		if menu.Type == MenuItem {
			m.root.correctCurCode()
			m.root.CurSubCode++
			menu.Code = fmt.Sprintf("%s%03d", m.root.Code, m.root.CurSubCode)
		} else {
			m.root.correctCurCode()
			m.root.CurFuncCode++
			menu.Code = fmt.Sprintf("%s_%02d", m.root.Code, m.root.CurFuncCode)
		}
	} else {
		if menu.Type == MenuItem {
			prtMenu := m.root.childrenMap[menu.ParentCode]
			prtMenu.correctCurCode()
			prtMenu.CurSubCode++
			menu.Code = fmt.Sprintf("%s%03d", prtMenu.Code, prtMenu.CurSubCode)
			if err = db.Menu().UpdateCurSubCode(prtMenu.Id, prtMenu.CurSubCode); err != nil {
				return err
			}
		} else {
			prtMenu := m.root.childrenMap[menu.ParentCode]
			prtMenu.correctCurCode()
			prtMenu.CurFuncCode++
			menu.Code = fmt.Sprintf("%s_%02d", prtMenu.Code, prtMenu.CurFuncCode)
			if err = db.Menu().UpdateCurFuncCode(prtMenu.Id, prtMenu.CurFuncCode); err != nil {
				return err
			}
		}
	}

	if menu.Id, err = db.Menu().Add(menu); err != nil {
		return err
	}

	reload()

	return nil
}

func (m *DataMenu) UpdateMenu(modJ model.Menu) error {
	var (
		err error
	)
	m.mux.Lock()
	defer m.mux.Unlock()

	if err = db.Menu().Update(modJ); err != nil {
		return err
	}

	reload()

	return nil
}

func (m *DataMenu) DeleteMenu(ids []int) error {
	var (
		err error
	)
	m.mux.Lock()
	defer m.mux.Unlock()

	if err = db.Menu().Delete(ids); err != nil {
		return err
	}

	reload()

	return nil
}

//构建选中菜单树
func (m *DataMenu) GetSelectMenuTree(menuCode string) []model.SelectMenuTree {
	//TODO
	return nil
}

//构建标记菜单树
func (m *DataMenu) GetCheckedMenuTree(menuCodes []string) []model.CheckedMenuTree {
	//TODO
	return nil
}

//获取权限菜单树
func (m *DataMenu) GetPowerMenuTree(menuCodes []string) []model.Power {
	//TODO
	return nil
}
