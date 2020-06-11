package api

import (
	"yumi/pkg/conf"
	"yumi/pkg/log"
	"yumi/pkg/net/ymhttp"
)

func DebugLog(c *ymhttp.Context) {
	if conf.IsDebug() {
		c.Next()
		return
	}

	log.Debug("req:", c.Request.URL.String())
	log.Debug("body:", c.Request.Body)
}

func Binding(c *ymhttp.Context) {

}
