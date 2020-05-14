package db

import (
	"errors"
	"fmt"
	"strings"
	"yumi/pkg/ecode"

	"yumi/consts"
	"yumi/pkg/external/dbc"
)

var (
	AcctRepeat = errors.New("用户名重复")
)

type DataAcct int

var u DataAcct

func Acct() DataAcct {
	return u
}

func (da DataAcct) Add(user, name, code, mobile, password, operator string) (id int64, err error) {
	sqlStr := `
		INSERT 
		INTO 
			power_accounts
			("user", "name", "code", "mobile", "password", "status", "register_time", "operator", "operate_time")
		VALUES
			(?, ?, ?, ?, ?, ?, sysdate(), ?, sysdate())`
	if id, err = dbc.Get().Insert(sqlStr,
		user, name, code, mobile, password, consts.AcctStatusEnable, operator); err != nil {
		if strings.Index(err.Error(), "UNIQUE KEY") != -1 {
			return 0, AcctRepeat
		}
		return 0, ecode.ServerErr(err)
	} else {
		return id, nil
	}
}

func (da DataAcct) Update(id int64, user, name, code, mobile, password, operator string) error {
	sets := make(map[string]interface{})
	sqlStr := `UPDATE power_accounts SET `
	if user != "" {
		sets["user"] = user
	}
	if code != "" {
		sets["code"] = code
	}
	if password != "" {
		sets["password"] = password
	}
	sets["name"] = name
	sets["mobile"] = mobile
	sets["operator"] = operator
	sets[`"operate_time"=sysdate()`] = nil

	setStr := ""
	setVals := []interface{}{}
	for k, v := range sets {
		if setStr == "" {
			setStr = fmt.Sprintf("%s=?", k)
			setVals = append(setVals, v)
		} else {
			setStr = fmt.Sprintf("%s, %s=?", setStr, k)
			setVals = append(setVals, v)
		}
	}

	sqlStr = fmt.Sprintf(`%s%s WHERE "id"=%d`, sqlStr, setStr, id)
	if _, err := dbc.Get().Exec(sqlStr, setVals...); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (da DataAcct) Delete(ids []int) error {
	idsstr := ""
	for i := range ids {
		if i == 0 {
			idsstr = fmt.Sprintf("'%d'", ids[i])
		} else {
			idsstr = fmt.Sprintf("%s,'%d'", idsstr, ids[i])
		}
	}

	sqlStr := `
			DELETE a, ar    
			FROM
				power_accounts AS a    
			LEFT JOIN
				power_acct_roles AS ar 
					ON ar.acct_code=a.code    
			WHERE
				FIND_IN_SET(a.id, ?)`
	if _, err := dbc.Get().Exec(sqlStr, idsstr); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (da DataAcct) Get(dest interface{}, id int64) error {
	if id == 0 {
		return nil
	} else {
		sqlStr := `
			SELECT 
				id, 
				ifnull("user", '') AS "user", 
				ifnull("name", '') AS "name", 
				ifnull("code", '') AS "code", 
				ifnull("mobile", '') AS "mobile", 
				ifnull("register_time", '') AS "registertime", 
				ifnull("status", '') AS "status", 
				ifnull("operator", '') AS "operator", 
				ifnull("operate_time", '') AS "operatetime" 
			FROM
			    power_accounts 
			WHERE
			    id=?`
		if err := dbc.Get().Get(dest, sqlStr, id); err != nil {
			return ecode.ServerErr(err)
		}
	}

	return nil
}

func (da DataAcct) Search(dest interface{}, user, name, code, mobile, operator string, page, line int) (
	total int, pageIndex int, pageCount int, err error) {
	cloumns := `id,
				isnull("user", '') AS "user", 
				isnull("user_id", '') AS "userid", 
				isnull("mobile", '') AS  "mobile", 
				isnull("comment", '') AS  "comment", 
				isnull("operator", '') AS  "operator", 
				isnull("register_time", '') AS  "registertime", 
				isnull("operate_time", '') AS "operatetime"`
	table := `power_accounts`
	conds := make(map[string]interface{})
	if user != "" {
		conds[`"user" LIKE ?`] = fmt.Sprintf(`%%%s%%`, user)
	}
	if name != "" {
		conds[`"name" LIKE ?`] = fmt.Sprintf(`%%%s%%`, name)
	}
	if mobile != "" {
		conds[`"mobile" LIKE ?`] = fmt.Sprintf(`%%%s%%`, mobile)
	}
	if operator != "" {
		conds[`"operator" LIKE ?`] = fmt.Sprintf(`%%%s%%`, operator)
	}
	where := ""
	whereVals := []interface{}{}
	for k, v := range conds {
		if where == "" {
			where = fmt.Sprintf("WHERE %s", k)
			whereVals = append(whereVals, v)
		} else {
			where = fmt.Sprintf("%s AND %s", where, k)
			whereVals = append(whereVals, v)
		}
	}
	order := `ORDER BY operate_time DESC`
	page, line = getDefaultPageLine(page, line)

	if total, pageIndex, pageCount, err = dbc.Get().PageSelect(dest, cloumns, table, where, order, page, line, whereVals...); err != nil {
		return 0, 0, 0, ecode.ServerErr(err)
	}

	return total, pageIndex, pageCount, nil
}

func (da DataAcct) SaveRolesOfAcct(acctCode string, roleCodes []string) error {
	tx, err := dbc.Get().Begin()
	if err != nil {
		return ecode.ServerErr(err)
	}
	defer tx.Rollback()

	sqlStr := `DELETE FROM power_acct_roles WHERE acct_code=?`
	if _, err = tx.Exec(sqlStr, acctCode); err != nil {
		return ecode.ServerErr(err)
	}

	stmt, err := tx.Prepare(`INSERT INTO power_acct_roles ("acct_code", "role_code") VALUES (?, ?)`)
	if err != nil {
		return ecode.ServerErr(err)
	}
	defer stmt.Close()
	for i := range roleCodes {
		if _, err = stmt.Exec(acctCode, roleCodes[i]); err != nil {
			return ecode.ServerErr(err)
		}
	}

	if err = tx.Commit(); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (da DataAcct) GetRolesOfAcct(dest interface{}, acctCode string) error {
	if acctCode == "" {
		sqlStr := `
			SELECT 
				id, 
				ifnull("name", '') AS "name", 
				ifnull("code", '') AS "code", 
				ifnull("operator", '') AS "operator", 
				ifnull("operate_time", '') AS "operatetime" 
			FROM 
			    power_roles 
			WHERE 
			    "status"=?`
		if err := dbc.Get().Select(dest, sqlStr, consts.RoleStatusEnable); err != nil {
			return ecode.ServerErr(err)
		}
	} else {
		sqlStr := `
			SELECT 
				r."id" AS "id", 
				ifnull(r."name", '') AS "name", 
				ifnull(r."code", '') AS "code", 
				ifnull(r."operator", '') AS "operator", 
				ifnull(r."operate_time", '') AS "operatetime" 
			FROM 
				power_roles AS r 
			LEFT JOIN 
				power_acct_roles AS ar 
				    ON ar."role_code"=r."code" 
			WHERE 
			    ar."acct_code"=?`
		if err := dbc.Get().Select(dest, sqlStr, acctCode); err != nil {
			return ecode.ServerErr(err)
		}
	}

	return nil
}

func (da DataAcct) SaveMenusOfAcct(acctCode string, menuCodes []string) error {
	tx, err := dbc.Get().Begin()
	if err != nil {
		return ecode.ServerErr(err)
	}
	defer tx.Rollback()

	sqlStr := `DELETE FROM power_acct_menus WHERE "acct_code"=?`
	if _, err := dbc.Get().Exec(sqlStr, acctCode); err != nil {
		return ecode.ServerErr(err)
	}

	sqlStr = `INSERT INTO power_acct_menus ("acct_code", "menu_code") VALUES (?, ?)`
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return ecode.ServerErr(err)
	}
	defer stmt.Close()
	for i := range menuCodes {
		if _, err := stmt.Exec(acctCode, menuCodes[i]); err != nil {
			return ecode.ServerErr(err)
		}
	}

	if err = tx.Commit(); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

func (da DataAcct) GetMenuCodesOfAcct(acctCode string) (menus []string, err error) {
	sqlStr := `SELECT "menu_code" FROM power_acct_menus WHERE acct_code=?`
	if err := dbc.Get().Select(&menus, sqlStr, acctCode); err != nil {
		return nil, ecode.ServerErr(err)
	}

	return menus, nil
}
