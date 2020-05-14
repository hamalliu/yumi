package conf

import (
	"github.com/BurntSushi/toml"

	"yumi/pkg/external/dbc/mysqlx"
)

const (
	EnvDebug   = "debug"
	EnvRelease = "release"
)

type Config struct {
	Server Server
	DB     mysqlx.Config
}

type Server struct {
	SysName      string //系统名称
	Addr         string //启动地址
	WriteTimeout int    //http写超时（second）
	ReadTimeout  int    //http读超时（second）
	Environment  string //运行环境
	StoragePath  string //附件路径
	MaxFileSize  int64  //附件最大限制
}

var conf Config

func init() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}
}

func Get() Server {
	return conf.Server
}

func GetDB() mysqlx.Config {
	return conf.DB
}
