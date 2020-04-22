package mssqlx

import (
	"time"

	"github.com/jmoiron/sqlx"
)

const dirverName = "mssql"

type Config struct {
	Dsn             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int64
}

type Model struct {
	*sqlx.DB
}

func New(conf Config) (*Model, error) {
	var (
		m   = new(Model)
		err error
	)

	if m.DB, err = sqlx.Connect(dirverName, conf.Dsn); err != nil {
		return nil, err
	}

	m.DB.SetMaxIdleConns(conf.MaxIdleConns)
	m.DB.SetMaxOpenConns(conf.MaxOpenConns)
	m.DB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Hour)

	return m, nil
}
