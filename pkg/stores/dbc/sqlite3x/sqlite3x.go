package sqlite3x

import (
	"github.com/jmoiron/sqlx"

	"yumi/pkg/stores/dbc"
)

const dirverName string = "sqlite3"

//Client sqlite3客户端
type Client struct {
	*sqlx.DB
}

//New 根据conf配置新建一个sqlite3客户端
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
