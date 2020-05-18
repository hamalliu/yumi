package mysqlx

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"yumi/pkg/conf"
	"yumi/pkg/log"
)

const dirverName = "mysql"

type Model struct {
	*sqlx.DB
	conf conf.DBConfig
}

func New(conf conf.DBConfig) (*Model, error) {
	var (
		m   = new(Model)
		err error
	)

	m.conf = conf
	if m.DB, err = sqlx.Connect(dirverName, conf.Dsn); err != nil {
		return nil, err
	}

	m.DB.SetMaxIdleConns(conf.MaxIdleConns)
	m.DB.SetMaxOpenConns(conf.MaxOpenConns)
	m.DB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Hour)

	//创建存储过程
	if f, err := os.Open("./page_select.sql"); err != nil {
		log.Error(err)
		return nil, err
	} else {
		if sqlb, err := ioutil.ReadAll(f); err != nil {
			log.Error(err)
			return nil, err
		} else {
			_, _ = m.Query(string(sqlb))
		}
	}

	return m, nil
}

func (m *Model) Insert(query string, args ...interface{}) (int64, error) {
	if res, err := m.Exec(query, args); err != nil {
		return 0, err
	} else if id, err := res.LastInsertId(); err != nil {
		return 0, err
	} else {
		return id, nil
	}

}

//分页查询
func (m *Model) PageSelect(dest interface{}, cloumns, table, where, order string, index, size int, args ...interface{}) (total, curIndex, curCount int, err error) {
	sqlStr := fmt.Sprintf(
		`
			declare 
				@total int; 
			declare 
				@cur_index int; 
			declare 
				@cur_count;
			call 
				page_select('%s', '%s', '%s', '%s', %d, %d, @total, @cur_index, @cur_count);
			select 
				@total, @cur_index, @cur_count;`,
		cloumns, table, where, order, index, size)
	if rows, err := m.Query(sqlStr, args...); err != nil {
		return 0, 0, 0, err
	} else {
		defer rows.Close()
		if err = sqlx.StructScan(rows, dest); err != nil {
			return 0, 0, 0, err
		} else {
			if rows.NextResultSet() {
				if rows.Next() {
					if err = rows.Scan(&total, &curIndex, &curCount); err != nil {
						return 0, 0, 0, err
					}
				}
			}
			return total, curIndex, curCount, nil
		}
	}
}

//软删除
func (m *Model) Delete(query string, args ...interface{}) error {
	var (
		sqls, ts = sqlDeleteToSelect(query)
		records  = make(map[string]string)
	)
	for i := range sqls {
		if rows, err := m.Query(sqls[i], args); err != nil {
			return err
		} else {
			records[ts[i]] = rowsToJson(rows)
		}
	}
	if recordsByte, err := json.Marshal(records); err != nil {
		return err
	} else {
		if err := m.insertDeleteRecords(query, string(recordsByte)); err != nil {
			return err
		}
	}

	if _, err := m.Exec(query, args); err != nil {
		return err
	}

	return nil
}

func (m *Model) insertDeleteRecords(query, records string) error {
	sqlStr := `
		INSERT
		INTO 
		    delete_records 
			("sql", "records", "operate_time") 
		VALUES 
			(?, ?, sysdate())`
	if _, err := m.Exec(sqlStr, query, records); err != nil {
		return err
	}

	return nil
}

func rowsToJson(rows *sql.Rows) string {
	var w bytes.Buffer
SCAN:
	valsMap := make(map[string]interface{})
	sqlx.MapScan(rows, valsMap)
	b, _ := json.Marshal(valsMap)
	w.Write(b)
	if rows.Next() {
		goto SCAN
	}

	return w.String()
}

func sqlDeleteToSelect(query string) ([]string, []string) {
	di := strings.Index(query, "FROM")
	head := query[:di]
	head = strings.Replace(head, "DELETE", "", -1)
	head = strings.ReplaceAll(head, " ", "")
	ts := strings.Split(head, ",")

	var sqls []string
	for i := range ts {
		sqlStr := fmt.Sprintf("SELECT %s.* %s", ts[i], query[di:])
		sqls = append(sqls, sqlStr)
	}

	return sqls, ts
}
