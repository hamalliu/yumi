package db

import (
	"fmt"

	"yumi/external/dbc"
	"yumi/model"
	"yumi/pkg/ecode"
)

type DataRollBackor int

func RollBackor() DataRollBackor {
	return 0
}

//新增更新记录
func (m DataRollBackor) AddUpdateRecord(table, request, requestbody, beforedata, userid, operator string) (int, error) {
	var (
		id  int64
		err error
	)
	sqlStr := `
			INSERT 
			INTO 
				back_update_records 
				("table", "request", "request_body", "before_data", "user_id", "operator", "operate_time") 
			VALUES 
				(?, ?, ?, ?, ?, ?, sysdate())`
	if id, err = dbc.Get().Insert(sqlStr, table, request, requestbody, beforedata, userid, operator); err != nil {
		return 0, ecode.ServerErr(err)
	}

	return int(id), nil
}

//删除更新记录
func (m DataRollBackor) DeleteUpdateRecord(ids []int) error {
	idsstr := ""
	for i := range ids {
		if i == 0 {
			idsstr = fmt.Sprintf("'%d'", ids[i])
		} else {
			idsstr = fmt.Sprintf("%s,'%d'", idsstr, ids[i])
		}
	}

	sqlStr := `DELETE FROM back_update_records WHERE FIND_IN_SET(id, ?)`
	if _, err := dbc.Get().Exec(sqlStr, idsstr); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

//查询更新记录
func (m DataRollBackor) SearchUpdateRecord(page, line int, startTime, endTime, userid string) ([]model.UpdateRecordsJson, int, int, int, error) {
	var (
		urjs      []model.UpdateRecordsJson
		total     int
		pageCount int
		pageIndex int

		where string
		conds []string

		err error
	)
	columns := `
		id, 
		isnull("table", '') AS "table", 
		isnull("request", '') AS "request", 
		isnull("request_body", '') AS "requestbody", 
		isnull("before_data", '') AS "beforedata", 
		isnull("user_id", '') AS "userid", 
		isnull("operate_time", '') AS "operatetime", 
		isnull("operator", '') AS "operator"`
	table := `back_update_records`
	if userid != "" {
		conds = append(conds, fmt.Sprintf("%s user_id = '%s'", where, userid))
	}
	if startTime != "" {
		conds = append(conds, fmt.Sprintf("%s operate_time > '%s'", where, startTime))
	}
	if endTime != "" {
		conds = append(conds, fmt.Sprintf("%s operate_time < '%s'", where, endTime))
	}
	for i := range conds {
		if where == "" {
			where = conds[i]
		} else {
			where = fmt.Sprintf("%s AND %s", where, conds[i])
		}
	}
	page, line = getDefaultPageLine(page, line)
	order := `"operate_time" desc`
	if total, pageCount, pageIndex, err = dbc.Get().PageSelect(&urjs, columns, table, where, order, page, line); err != nil {
		return nil, 0, 0, 0, ecode.ServerErr(err)
	}

	return urjs, total, pageCount, pageIndex, nil
}

//获取删除记录
func (m DataRollBackor) GetUpdateRecords(ids []int) ([]model.UpdateRecordsJson, error) {
	var (
		udjs []model.UpdateRecordsJson
		err  error
	)
	idsstr := ""
	for i := range ids {
		if i == 0 {
			idsstr = fmt.Sprintf("'%d'", ids[i])
		} else {
			idsstr = fmt.Sprintf("%s,'%d'", idsstr, ids[i])
		}
	}

	sqlStr := `
			SELECT 
				id, 
				ifnull("table", '') AS "table", 
				ifnull("request", '') AS "request", 
				ifnull("request_body", '') AS "requestbody", 
				ifnull("before_data", '') AS "beforedata", 
				ifnull("user_id", '') AS "userid", 
				ifnull("operate_time", '') AS "operatetime", 
				ifnull("operator", '') AS "operator" 
			FROM 
			    back_update_records 
			WHERE 
			    FIND_IN_SET(id, ?)`
	if err = dbc.Get().Select(&udjs, sqlStr, idsstr); err != nil {
		return nil, ecode.ServerErr(err)
	}

	return udjs, nil
}

//新增删除记录
func (m DataRollBackor) AddDeleteRecord(table, request, requestbody, beforedata, userid, operator string) (int, error) {
	var (
		id  int64
		err error
	)

	sqlStr := `
			INSERT 
			INTO 
				back_delete_records
				("table", "request", "request_body", "before_data", "user_id", "operator", "operate_time")
			VALUES 
				(?, ?, ?, ?, ?, ?, sysdate())`
	if id, err = dbc.Get().Insert(sqlStr, table, request, requestbody, beforedata, userid, operator); err != nil {
		return 0, ecode.ServerErr(err)
	}

	return int(id), nil
}

//删除删除记录
func (m DataRollBackor) DeleteDeleteRecord(ids []int) error {
	idsstr := ""
	for i := range ids {
		if i == 0 {
			idsstr = fmt.Sprintf("'%d'", ids[i])
		} else {
			idsstr = fmt.Sprintf("%s,'%d'", idsstr, ids[i])
		}
	}

	sqlStr := `DELETE FROM back_delete_records WHERE FIND_IN_SET(id, ?)`
	if _, err := dbc.Get().Exec(sqlStr, idsstr); err != nil {
		return ecode.ServerErr(err)
	}

	return nil
}

//查询删除记录
func (m DataRollBackor) SearchDeleteRecord(page, line int, startTime, endTime, userid string) ([]model.DeleteRecordsJson, int, int, int, error) {
	var (
		drjs      []model.DeleteRecordsJson
		total     int
		pageCount int
		pageIndex int

		where string
		conds []string

		err error
	)
	colmnus := `
		id, 
		isnull("table", '') AS "table", 
		isnull("request", '') AS "request", 
		isnull("request_body", '') AS "requestbody", 
		isnull("before_data", '') AS "beforedata", 
		isnull("user_id", '') AS "userid", 
		isnull("operate_time", '') AS "operatetime", 
		isnull("operator", '') AS "operator", `
	table := `back_delete_records`
	if userid != "" {
		conds = append(conds, fmt.Sprintf("%s user_id = '%s'", where, userid))
	}
	if startTime != "" {
		conds = append(conds, fmt.Sprintf("%s operate_time > '%s'", where, startTime))
	}
	if endTime != "" {
		conds = append(conds, fmt.Sprintf("%s operate_time < '%s'", where, endTime))
	}
	for i := range conds {
		if where == "" {
			where = conds[i]
		} else {
			where = fmt.Sprintf("%s AND %s", where, conds[i])
		}
	}
	order := `"operate_time" desc`
	page, line = getDefaultPageLine(page, line)
	if total, pageCount, pageIndex, err = dbc.Get().PageSelect(&drjs, colmnus, table, where, order, page, line); err != nil {
		return nil, 0, 0, 0, ecode.ServerErr(err)
	}

	return drjs, total, pageCount, pageIndex, nil
}

//获取删除记录
func (m DataRollBackor) GetDeleteRecords(ids []int) ([]model.DeleteRecordsJson, error) {
	var (
		udjs []model.DeleteRecordsJson
		err  error
	)
	idsstr := ""
	for i := range ids {
		if i == 0 {
			idsstr = fmt.Sprintf("'%d'", ids[i])
		} else {
			idsstr = fmt.Sprintf("%s,'%d'", idsstr, ids[i])
		}
	}

	sqlStr := `
			SELECT 
				id, 
				ifnull("table", '') AS "table", 
				ifnull("request", '') AS "request", 
				ifnull("request_body", '') AS "requestbody", 
				ifnull("before_data", '') AS "beforedata", 
				ifnull("user_id", '') AS "userid", 
				ifnull("operate_time", '') AS "operatetime", 
				ifnull("operator", '') AS "operator" 
			FROM 
			    back_delete_records 
			WHERE 
			    FIND_IN_SET(id, ?)`
	if err = dbc.Get().Select(&udjs, sqlStr, idsstr); err != nil {
		return nil, ecode.ServerErr(err)
	}

	return udjs, nil
}
