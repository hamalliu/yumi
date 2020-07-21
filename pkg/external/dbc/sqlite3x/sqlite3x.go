package sqlite3x

import (
	"yumi/pkg/conf"

	"github.com/jmoiron/sqlx"
)

const dirverName string = "sqlite3"

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
