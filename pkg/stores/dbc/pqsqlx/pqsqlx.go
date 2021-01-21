package pqsqlx

import (
	"yumi/conf"

	"github.com/jmoiron/sqlx"
)

const dirverName = "postgres"

//Client postgres 客户端
type Client struct {
	*sqlx.DB
}

//New 根据conf配置新建数据库客户端
func New(conf conf.DB) (*Client, error) {
	var (
		m   = new(Client)
		err error
	)

	if m.DB, err = sqlx.Connect(dirverName, conf.Dsn); err != nil {
		return nil, err
	}

	m.DB.SetMaxIdleConns(conf.MaxIdleConns)
	m.DB.SetMaxOpenConns(conf.MaxOpenConns)
	m.DB.SetConnMaxLifetime(conf.ConnMaxLifetime.Duration())

	return m, nil
}
