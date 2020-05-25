package handler_conf

type ConfigData map[string]Config

type Config struct {
	Encrypt bool
}

type ConfigInf interface {
	Load() ConfigData
}

var _ci ConfigInf

func Register(ci ConfigInf) {
	_ci = ci
}

var _cd ConfigData

func Init() {
	_cd = _ci.Load()
}

func Get(code string) Config {
	return _cd[code]
}
