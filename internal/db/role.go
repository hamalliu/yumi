package db

import (
	"errors"
	"fmt"
	"strings"

	"yumi/consts"
	"yumi/model"
	"yumi/pkg/ecode"
	"yumi/pkg/external/dbc"
)

var (
	RoleRepeat = errors.New("角色名重复")
)

type DataRole int

var r DataRole

func Role() DataRole {
	return r
}

func (m DataRole) Add(name, code, operator string) (id int64, err error) {
	sqlStr := `
		INSERT 
		INTO
    		power_roles
    		("name", "code", "status", "operator", "operate_time")     
		VALUES
    		(?, ?, ?, ?, getdate())`
	if id, err = dbc.Get().Insert(sqlStr, name, code, consts.RoleStatusEnable, operator); err != nil {
		if strings.Index(err.Error(), "UNIQUE KEY") != -1 {
			return 0, RoleRepeat
		}

		return 0, ecode.ServerErr(err)
	}

	return id, nil
}

func (m DataRole) Update(id int, name, code, status, operator string) error {
	if name == "" {
		return nil
	}

	sqlStr := `
		UPDATE
    		power_roles 
		SET
    		"name"=?,
			"code"=?,
			"status"=?,
			"operator"=?,
			"operate_time"=getdate()      
		WHERE
			id=?`
	if _, err := dbc.Get().Exec(sqlStr, name, code, status, operator, id); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (m DataRole) Delete(ids []int) error {
	idsstr := ""
	for i := range ids {
		if i == 0 {
			idsstr = fmt.Sprintf("'%d'", ids[i])
		} else {
			idsstr = fmt.Sprintf("%s,'%d'", idsstr, ids[i])
		}
	}

	sqlStr := `
		DELETE ur, r 
		FROM
			power_acct_roles AS ur      
		LEFT JOIN
			power_roles r 
				ON r.code=ur.role_code 
		WHERE
			FIND_IN_SET (r.id, ?)`
	if _, err := dbc.Get().Exec(sqlStr, idsstr); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (m DataRole) Get(id int) (role model.Role, err error) {
	sqlStr := `
		SELECT
			id,
			ifnull("name", '') AS "name", 
			ifnull("code", '') AS "code",
			ifnull("status", '') AS "status", 
			ifnull("operator", '') AS "operator",
			ifnull("operate_time", '') AS "operatetime"      
		FROM
			power_roles   
		WHERE
			id=?`
	if err = dbc.Get().Get(&role, sqlStr, id); err != nil {
		return role, ecode.ServerErr(err)
	}

	return role, nil
}

func (m DataRole) Search(page, line int, name, code, status string) (
	roles []model.Role, total int, pageIndex int, pageCount int, err error) {
	cloumns := `
				id, 
				isnull("name", '') AS "name",
				isnull("code", '') AS "code",
				isnull("status", '') AS "status",
				isnull("operator", '') AS "operator",
				isnull("operate_time", '') AS "operatetime"`
	table := `power_roles`
	order := `"operate_time" desc`

	conds := []string{}
	if status != "" {
		conds = append(conds, fmt.Sprintf("status='%s'", status))
	}
	if name != "" {
		conds = append(conds, fmt.Sprintf("name LIKE '%%%s%%'", name))
	}
	if code != "" {
		conds = append(conds, fmt.Sprintf("code LIKE '%%%s%%'", code))
	}
	where := ""
	for i := range conds {
		if where == "" {
			where = fmt.Sprintf("WHERE %s", conds[i])
		} else {
			where = fmt.Sprintf("%s AND %s", where, conds[i])
		}
	}
	page, line = getDefaultPageLine(page, line)
	if total, pageIndex, pageCount, err = dbc.Get().PageSelect(&roles, cloumns, table, where, order, page, line); err != nil {
		return nil, 0, 0, 0, ecode.ServerErr(err)
	}

	return roles, total, pageIndex, pageCount, nil
}

func (m DataRole) SaveUsersOfRole(roleCode string, acctCodes []string) error {
	tx, err := dbc.Get().Begin()
	if err != nil {
		return ecode.ServerErr(err)
	}
	defer tx.Rollback()

	sqlStr := `DELETE FROM power_acct_roles WHERE role_code=?`
	if _, err = tx.Exec(sqlStr, roleCode); err != nil {
		return ecode.ServerErr(err)
	}

	sqlStr = `INSERT INTO power_acct_roles ("acct_code", "role_code") VALUES (?, ?)`
	if stmt, err := tx.Prepare(sqlStr); err != nil {
		return ecode.ServerErr(err)
	} else {
		defer stmt.Close()
		for i := range acctCodes {
			if _, err = stmt.Exec(acctCodes[i], roleCode); err != nil {
				return ecode.ServerErr(err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (m DataRole) GetUsersOfRole(roleCode string) (accts []model.Account, err error) {
	if roleCode == "" {
		sqlStr := `
		SELECT 
       		"id", 
       		ifnull("user", '') AS "user", 
       		ifnull("name", '') AS "name", 
       		ifnull("code", '') AS "code", 
			ifnull("mobile", '') AS "mobile", 
			ifnull("register_time", '') AS "register_time",
			ifnull("operator", '') AS "operator", 
			ifnull("operate_time", '') AS "operate_time" 
		FROM 
			power_accounts 
		WHERE 
			"status"=?`
		if err = dbc.Get().Select(&accts, sqlStr, consts.AcctStatusEnable); err != nil {
			return nil, ecode.ServerErr(err)
		}
	} else {
		sqlStr := `
		SELECT 
       		acct."id" AS "id", 
       		ifnull(acct."user", '') AS "user", 
       		ifnull(acct."name", '') AS "name", 
       		ifnull(acct."code", '') AS "code", 
			ifnull(acct."mobile", '') AS "mobile", 
			ifnull(acct."register_time", '') AS "register_time",
			ifnull(acct."operator", '') AS "operator", 
			ifnull(acct."operate_time", '') AS "operate_time" 
		FROM 
		    power_accounts AS acct 
		LEFT JOIN 
			power_acct_roles AS ar 
		    	ON ar.acct_code = acct.code 
		WHERE 
		    ar."role_code"=? 
		  	AND acct."status"=?`
		if err = dbc.Get().Select(&accts, sqlStr, roleCode, consts.AcctStatusEnable); err != nil {
			return nil, ecode.ServerErr(err)
		}
	}

	return accts, nil
}

func (m DataRole) SaveMenusOfRole(roleCode string, menusCodes []string) error {
	tx, err := dbc.Get().Begin()
	if err != nil {
		return ecode.ServerErr(err)
	}
	defer tx.Rollback()

	sqlStr := `DELETE FROM power_role_menus WHERE role_code=?`
	if _, err := tx.Exec(sqlStr, roleCode); err != nil {
		return ecode.ServerErr(err)
	}

	sqlStr = `INSERT INTO power_role_menus ("role_code", "menu_code") VALUES (?, ?)`
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return ecode.ServerErr(err)
	}

	for i := range menusCodes {
		if _, err := stmt.Exec(roleCode, menusCodes[i]); err != nil {
			return ecode.ServerErr(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (m DataRole) GetMenuCodesOfRole(roleCode string) (menus []string, err error) {
	sqlStr := `SELECT ifnull("menu_code", '') AS "code" FROM power_role_menus WHERE role_code=?`
	if err = dbc.Get().Select(&menus, sqlStr, roleCode); err != nil {
		return nil, ecode.ServerErr(err)
	}

	return menus, nil
}
