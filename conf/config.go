package conf

import (
	"github.com/BurntSushi/toml"

	"yumi/pkg/types"
)

// Config app config
type Config struct {
	Deploy     Deploy
	Server     Server
	Media      Media
	DB         DB
	Mongo      Mongo
	CORS       CORS
	OnlyOffice OnlyOffice
}

// Server 服务器配置
type Server struct {
	Addr         string             // 启动地址
	WriteTimeout types.TimeDuration // http写超时
	ReadTimeout  types.TimeDuration // http读超时
}

// Media 媒体配置
type Media struct {
	StoragePath                string          // 附件路径
	MultipleFileUploadsMaxSize types.SpaceSize // 多媒体上传最大限制
	SingleFileUploadsMaxSize   types.SpaceSize // 单媒体上传最大限制
}

// CORS CORS配置
type CORS struct {
	AllowedOrigins []string           // 允许的头
	MaxAge         types.TimeDuration // 最大持续时间
}

var conf Config

// Load 加载配置
func Load() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}
}

// Get 获取配置
func Get() Config {
	return conf
}

// IsDebug 该程序是调试模式
func IsDebug() bool {
	return conf.Deploy.Env == DeployEnvDev
}
