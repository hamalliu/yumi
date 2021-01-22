package conf

import (
	"yumi/pkg/stores/dbc"
	"yumi/pkg/types"
)

// DB 数据库配置
type DB struct {
	Dsn             string             // 数据来源名称
	DBName          string             // 数据库名称
	MaxOpenConns    int                // 连接池最多打开连接数
	MaxIdleConns    int                // 连接池最多空闲连接数
	ConnMaxLifetime types.TimeDuration // 连接最长寿命
}

// Options ...
func (c *DB) Options() []dbc.ClientOption {
	opts := []dbc.ClientOption{}
	opts = append(opts, dbc.SetMaxIdleConns(c.MaxIdleConns))
	opts = append(opts, dbc.SetMaxIdleConns(c.MaxOpenConns))
	opts = append(opts, dbc.SetConnMaxLifetime(c.ConnMaxLifetime.Duration()))

	return opts
}
