package log

import (
	"os"

	"github.com/op/go-logging"

	"yumi/conf"
)

var infolog = logging.MustGetLogger("info")
var errorlog = logging.MustGetLogger("error")
var criticallog = logging.MustGetLogger("critical")
var log = logging.MustGetLogger("log")

func init() {
	var format = logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05.000} %{level:.4s} %{shortfile} %{message}`,
	)
	logging.SetFormatter(format)

	stdBackend := logging.NewLogBackend(os.Stdout, "", 0)
	stdBackend.Color = true

	infoBackend := logging.NewLogBackend(New("[INFO]", 2<<26, true), "", 0)
	lvlInfoBackend := logging.AddModuleLevel(infoBackend)
	lvlInfoBackend.SetLevel(logging.ERROR, "")

	errBackend := logging.NewLogBackend(New("[ERROR]", 2<<26, true), "", 0)
	lvlErrBackend := logging.AddModuleLevel(errBackend)
	lvlErrBackend.SetLevel(logging.ERROR, "")

	criticalBackend := logging.NewLogBackend(New("[CRITICAL]", 2<<26, true), "", 0)
	lvlCriticalBackend := logging.AddModuleLevel(criticalBackend)
	lvlCriticalBackend.SetLevel(logging.CRITICAL, "")

	logging.SetBackend(stdBackend)
	infolog.SetBackend(lvlInfoBackend)
	errorlog.SetBackend(lvlErrBackend)
	criticallog.SetBackend(lvlCriticalBackend)
}

func Critical(args ...interface{}) {
	criticallog.ExtraCalldepth = 1
	criticallog.Critical(args)
	if conf.Get().Environment == conf.EnvDebug {
		log.ExtraCalldepth = 1
		log.Debug(args)
	}
}

func Error(args ...interface{}) {
	errorlog.ExtraCalldepth = 1
	errorlog.Error(args)
	if conf.Get().Environment == conf.EnvDebug {
		log.ExtraCalldepth = 1
		log.Debug(args)
	}
}

func Info(args ...interface{}) {
	infolog.ExtraCalldepth = 1
	infolog.Error(args)
	if conf.Get().Environment == conf.EnvDebug {
		log.ExtraCalldepth = 1
		log.Debug(args)
	}
}

func Debug(args ...interface{}) {
	log.ExtraCalldepth = 1
	log.Debug(args)
}

func Critical2(args ...interface{}) {
	criticallog.ExtraCalldepth = 2
	criticallog.Critical(args)
	if conf.Get().Environment == conf.EnvDebug {
		log.ExtraCalldepth = 2
		log.Debug(args)
	}
}

func Error2(args ...interface{}) {
	errorlog.ExtraCalldepth = 2
	errorlog.Error(args)
	if conf.Get().Environment == conf.EnvDebug {
		log.ExtraCalldepth = 2
		log.Debug(args)
	}
}

func Info2(args ...interface{}) {
	infolog.ExtraCalldepth = 2
	infolog.Error(args)
	if conf.Get().Environment == conf.EnvDebug {
		log.ExtraCalldepth = 2
		log.Debug(args)
	}
}
