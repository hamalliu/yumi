package dbc

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// ClientOption ...
type ClientOption struct {
	F func(*ClientOptions)
}

// ClientOptions ...
type ClientOptions struct {
	DB *sqlx.DB
}

// SetMaxIdleConns ...
func SetMaxIdleConns(n int) ClientOption {
	return ClientOption{
		F: func(co *ClientOptions) {
			if n != 0 {
				co.DB.SetMaxIdleConns(n)
			}
		},
	}
}

// SetMaxOpenConns ...
func SetMaxOpenConns(n int) ClientOption {
	return ClientOption{
		F: func(co *ClientOptions) {
			if n != 0 {
				co.DB.SetMaxOpenConns(n)
			}
		},
	}
}

// SetConnMaxLifetime ...
func SetConnMaxLifetime(d time.Duration) ClientOption {
	return ClientOption{
		F: func(co *ClientOptions) {
			if d != 0 {
				co.DB.SetConnMaxLifetime(d)
			}
		},
	}
}
