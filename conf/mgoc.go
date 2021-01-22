package conf

import (
	"go.mongodb.org/mongo-driver/mongo/options"

	"yumi/pkg/types"
)

// Mongo mongodb 数据库配置
type Mongo struct {
	Dsn               string
	DBName            string
	MaxPoolSize       uint64             // 连接池最多打开连接数
	MinPoolSize       uint64             // 连接池最少打开连接数
	MaxConnIdleTime   types.TimeDuration // 空闲连接断开时间
	ConnectTimeout    types.TimeDuration // 连接超时时间
	HeartbeatInterval types.TimeDuration // 心跳间隔时间
}

// Options ...
func (c *Mongo) Options() []*options.ClientOptions {
	opts := []*options.ClientOptions{}
	if c.ConnectTimeout != 0 {
		opts = append(opts, options.Client().SetConnectTimeout(c.ConnectTimeout.Duration()))
	}
	if c.HeartbeatInterval != 0 {
		opts = append(opts, options.Client().SetHeartbeatInterval(c.HeartbeatInterval.Duration()))
	}
	if c.MaxConnIdleTime != 0 {
		opts = append(opts, options.Client().SetMaxConnIdleTime(c.MaxConnIdleTime.Duration()))
	}
	if c.MaxPoolSize != 0 {
		opts = append(opts, options.Client().SetMaxPoolSize(c.MaxPoolSize))
	}
	if c.MinPoolSize != 0 {
		opts = append(opts, options.Client().SetMinPoolSize(c.MinPoolSize))
	}

	return opts
}
