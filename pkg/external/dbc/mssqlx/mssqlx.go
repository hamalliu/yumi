package mssqlx

import (
	"github.com/jmoiron/sqlx"

	"yumi/pkg/conf"
)

const dirverName = "mssql"

type Model struct {
	*sqlx.DB
}

func New(conf conf.DB) (*Model, error) {
	var (
		m   = new(Model)
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
