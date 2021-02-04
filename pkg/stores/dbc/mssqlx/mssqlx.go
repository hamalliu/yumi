package mssqlx

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"

	"yumi/pkg/stores/dbc"
)

const dirverName = "mssql"

//Client mssql 客户端
type Client struct {
	*sqlx.DB
}

//New 新建一个mssql客户端
func New(dsn string, options ...dbc.ClientOption) (*Client, error) {
	var (
		m   = new(Client)
		err error
	)

	if m.DB, err = sqlx.Connect(dirverName, dsn); err != nil {
		return nil, err
	}

	opts := &dbc.ClientOptions{DB: m.DB}
	for _, option := range options {
		option.F(opts)
	}

	return m, nil
}
