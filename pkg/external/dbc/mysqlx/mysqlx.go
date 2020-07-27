package mysqlx

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// 导入驱动包
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"yumi/pkg/conf"
	"yumi/pkg/log"
)

const dirverName = "mysql"

//Client mysql 客户端
type Client struct {
	*sqlx.DB
	conf conf.DB
}

//New 新建一个 mysql 客户端
func New(conf conf.DB) (*Client, error) {
	var (
		m   = new(Client)
		err error
	)

	m.conf = conf
	if m.DB, err = sqlx.Connect(dirverName, conf.Dsn); err != nil {
		return nil, err
	}

	m.DB.SetMaxIdleConns(conf.MaxIdleConns)
	m.DB.SetMaxOpenConns(conf.MaxOpenConns)
	m.DB.SetConnMaxLifetime(conf.ConnMaxLifetime.Duration())

	//创建存储过程
	f, err := os.Open("./page_select.sql")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	sqlb, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	_, _ = m.Query(string(sqlb))

	return m, nil
}

//Insert exec insert sql
func (m *Client) Insert(query string, args ...interface{}) (int64, error) {
	if res, err := m.Exec(query, args); err != nil {
		return 0, err
	} else if id, err := res.LastInsertId(); err != nil {
		return 0, err
	} else {
		return id, nil
	}

}

//PageSelect 分页查询
func (m *Client) PageSelect(dest interface{}, cloumns, table, where, order string, index, size int, args ...interface{}) (total, curIndex, curCount int, err error) {
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
	rows, err := m.Query(sqlStr, args...)
	if err != nil {
		return 0, 0, 0, err
	}
	defer rows.Close()
	if err = sqlx.StructScan(rows, dest); err != nil {
		return 0, 0, 0, err
	}
	if rows.NextResultSet() {
		if rows.Next() {
			if err = rows.Scan(&total, &curIndex, &curCount); err != nil {
				return 0, 0, 0, err
			}
		}
	}
	return total, curIndex, curCount, nil
}

//Delete 软删除
func (m *Client) Delete(query string, args ...interface{}) error {
	var (
		sqls, ts = sqlDeleteToSelect(query)
		records  = make(map[string]string)
	)
	for i := range sqls {
		rows, err := m.Query(sqls[i], args)
		if err != nil {
			return err
		}
		records[ts[i]] = rowsToJSON(rows)
	}
	recordsByte, err := json.Marshal(records)
	if err != nil {
		return err
	}
	if err := m.insertDeleteRecords(query, string(recordsByte)); err != nil {
		return err
	}

	if _, err := m.Exec(query, args); err != nil {
		return err
	}

	return nil
}

func (m *Client) insertDeleteRecords(query, records string) error {
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

func rowsToJSON(rows *sql.Rows) string {
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
