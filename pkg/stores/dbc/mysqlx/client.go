package mysqlx

import (
	"context"
	"database/sql"
)

type CURD interface {
	Insert(query string, args ...interface{}) (int64, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Delete(query string, args ...interface{}) error
	PageSelect(dest interface{}, cloumns, table, where, order string, index, size int, args ...interface{}) (
		total, curIndex, curCount int, err error)
}
