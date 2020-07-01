package api

import (
	"yumi/pkg/conf"
	"yumi/pkg/log"
	"yumi/pkg/net/gin"
)

func DebugLog(c *gin.Context) {
	if conf.IsDebug() {
		c.Next()
		return
	}

	log.Debug("req:", c.Request.URL.String())
	log.Debug("body:", c.Request.Body)
}

func Binding(c *gin.Context) {

}
