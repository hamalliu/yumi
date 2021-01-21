package api

import (
	"yumi/conf"
	"yumi/gin"
	"yumi/pkg/log"
)

//DebugLog ...
func DebugLog(c *gin.Context) {
	if conf.IsDebug() {
		c.Next()
		return
	}

	log.Debug("req:", c.Request.URL.String())
	log.Debug("body:", c.Request.Body)
}
