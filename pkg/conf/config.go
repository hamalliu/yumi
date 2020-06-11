package conf

import (
	"github.com/BurntSushi/toml"
)

const (
	envDebug   = "debug"
	envRelease = "release"
)

type Config struct {
	Program Program
	Server  Server
	Media   Media
	DB      DB
	CORS    CORS
}

type Program struct {
	SysName     string //系统名称
	Environment string //运行环境
}

type Server struct {
	Addr         string       //启动地址
	WriteTimeout TimeDuration //http写超时
	ReadTimeout  TimeDuration //http读超时
}

type Media struct {
	StoragePath                string    //附件路径
	MultipleFileUploadsMaxSize SpaceSize //多媒体上传最大限制
	SingleFileUploadsMaxSize   SpaceSize //单媒体上传最大限制
}

type DB struct {
	Dsn             string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime TimeDuration
}

type CORS struct {
	AllowedOrigins []string
	MaxAge         TimeDuration
}

var conf Config

func Load() {
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}
}

func Get() Config {
	return conf
}

func IsDebug() bool {
	return conf.Program.Environment == envDebug
}
