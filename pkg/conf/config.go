package conf

import (
	"github.com/BurntSushi/toml"
)

const (
	EnvDebug   = "debug"
	EnvRelease = "release"
)

type Config struct {
	Server Server
	DB     DBConfig
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

type DBConfig struct {
	Dsn             string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int64
}

var conf Config

func Load() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}
}

func Get() Server {
	return conf.Server
}

func GetDB() DBConfig {
	return conf.DB
}
