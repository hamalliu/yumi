package db

import (
	"fmt"

	"yumi/consts"
	"yumi/model"
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

type DataMenu int

var menu DataMenu

func Menu() DataMenu {
	return menu
}

func (DataMenu) Add(menu model.Menu) (int64, error) {
	sqlStr := `
			INSERT 
			INTO
				power_menus
				("parent_name", "parent_code", "name", "code", "route", "params", "type", "display_order", "status", "cur_sub_code", "cur_func_code", "operator", "operate_time")    
			VALUES
				(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, sysdate())`
	if id, err := dbc.Get().Insert(sqlStr,
		menu.ParentName, menu.ParentCode, menu.Name, menu.Code, menu.Route,
		menu.Params, menu.Type, menu.DisplayOrder, menu.Status, menu.CurSubCode,
		menu.CurFuncCode, menu.Operator); err != nil {
		return 0, ecode.ServerErr(err)
	} else {
		return id, nil
	}
}

func (DataMenu) Delete(modids []int) error {
	idsstr := ""
	for i := range modids {
		if i == 0 {
			idsstr = fmt.Sprintf("'%d'", modids[i])
		} else {
			idsstr = fmt.Sprintf(",'%d'", modids[i])
		}
	}

	sqlStr := `DELETE FROM power_menus WHERE FIND_IN_SET(id, ?)`
	if _, err := dbc.Get().Exec(sqlStr, idsstr); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (DataMenu) Update(menu model.Menu) error {
	sqlStr := `
			UPDATE
				power_menus 
			SET
				"parent_name"=?,
				"parent_code"=?,
				"name"=?,
				"code"=?,
				"route"=?,
				"params"=?,
				"type"=?,
				"display_order"=?,
				"status"=?,
				"cur_sub_code"=?,
				"cur_func_code"=?,
				"operator"=?,
				"operate_time"=sysdate()    
			WHERE
				"id"=?`
	if _, err := dbc.Get().Exec(sqlStr,
		menu.ParentName, menu.ParentCode, menu.Name, menu.Code, menu.Route,
		menu.Params, menu.Type, menu.DisplayOrder, menu.Status, menu.CurSubCode,
		menu.CurFuncCode, menu.Operator, menu.Id); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (DataMenu) Get(id int64) (model.Menu, error) {
	menu := model.Menu{}

	if id == 0 {
		return menu, nil
	}

	sqlStr := `
		SELECT 
			id, 
			ifnull("parent_name", 0) AS "parentname",
			ifnull("parent_code", '') AS "parentcode", 
			ifnull("name", '') AS "name", 
			ifnull("code", '') AS "code", 
			ifnull("route", '') AS "route", 
			ifnull("params", '') AS "params", 
			ifnull("type", '') AS "type", 
			ifnull("display_order", '') AS "displayorder", 
			ifnull("status", 0) AS "status", 
			ifnull("cur_sub_code", '') AS "cursubcode", 
			ifnull("cur_func_code", '') AS "curfunccode", 
			ifnull("operator", '') AS "operator", 
			ifnull("operate_time", '') AS "operatetime" 
		FROM 
		     power_menus 
		WHERE
		     "id"=?`
	if err := dbc.Get().Get(&menu, sqlStr, id); err != nil {
		return menu, ecode.ServerErr(err)
	}

	return menu, nil
}

func (DataMenu) GetEnableMenus() ([]model.Menu, error) {
	menus := []model.Menu{}

	sqlStr := `
		SELECT 
			id, 
			ifnull("parent_name", 0) AS "parentname",
			ifnull("parent_code", '') AS "parentcode", 
			ifnull("name", '') AS "name", 
			ifnull("code", '') AS "code", 
			ifnull("route", '') AS "route", 
			ifnull("params", '') AS "params", 
			ifnull("type", '') AS "type", 
			ifnull("display_order", '') AS "displayorder", 
			ifnull("status", 0) AS "status", 
			ifnull("cur_sub_code", '') AS "cursubcode", 
			ifnull("cur_func_code", '') AS "curfunccode", 
			ifnull("operator", '') AS "operator", 
			ifnull("operate_time", '') AS "operatetime" 
		FROM
		    power_menus 
		WHERE
		    status=? 
		ORDER BY
		    "display_order"`
	if err := dbc.Get().Select(&menus, sqlStr, consts.MenuStatusEnable); err != nil {
		return nil, ecode.ServerErr(err)
	}

	return menus, nil
}

func (DataMenu) GetAllMenus() ([]model.Menu, error) {
	menus := []model.Menu{}

	sqlStr := `
		SELECT 
			id, 
			ifnull("parent_name", 0) AS "parentname",
			ifnull("parent_code", '') AS "parentcode", 
			ifnull("name", '') AS "name", 
			ifnull("code", '') AS "code", 
			ifnull("route", '') AS "route", 
			ifnull("params", '') AS "params", 
			ifnull("type", '') AS "type", 
			ifnull("display_order", '') AS "displayorder", 
			ifnull("status", 0) AS "status", 
			ifnull("cur_sub_code", '') AS "cursubcode", 
			ifnull("cur_func_code", '') AS "curfunccode", 
			ifnull("operator", '') AS "operator", 
			ifnull("operate_time", '') AS "operatetime" 
		FROM
		    power_menus 
		ORDER BY
		    "display_order"`
	if err := dbc.Get().Select(&menus, sqlStr); err != nil {
		return nil, ecode.ServerErr(err)
	}

	return menus, nil
}

func (DataMenu) UpdateCurSubCode(id int64, curSubCode uint) error {
	sqlStr := `UPDATE power_menus SET cur_sub_code=? WHERE id=?`
	if _, err := dbc.Get().Exec(sqlStr, curSubCode, id); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (DataMenu) UpdateCurFuncCode(id int64, curFuncCode uint) error {
	sqlStr := `UPDATE power_menus SET cur_func_code=? WHERE id=?`
	if _, err := dbc.Get().Exec(sqlStr, curFuncCode, id); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}
