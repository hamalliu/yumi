package conf

import (
	"github.com/BurntSushi/toml"

	"yumi/pkg/types"
)

// Config app config
type Config struct {
	// 运维层配置
	HostEnv    HostEnv
	HttpServer HTTPServer
	DB         DB
	Mongo      Mongo
	Log        Log
	Casbin     Casbin

	// 业务层配置
	Media      Media
	OnlyOffice OnlyOffice
	Trade      Trade
}

// HTTPServer 服务器配置
type HTTPServer struct {
	Addr           string             // 启动地址
	WriteTimeout   types.TimeDuration // http写超时
	ReadTimeout    types.TimeDuration // http读超时
	HandlerTimeout types.TimeDuration // 请求处理超时

	// CORS 配置
	CORSAllowedOrigins []string           // 允许的头
	CORSMaxAge         types.TimeDuration // 最大持续时间
}

type Casbin struct {
	ModelFile string
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
	return conf.HostEnv.Environment == DeployEnvDev
}
