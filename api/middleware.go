package api

import (
	"yumi/internal/session"
	"yumi/pkg/conf"
	"yumi/pkg/log"
	"yumi/pkg/net/gin"
	"yumi/pkg/net/gin/header"
	"yumi/pkg/valuer"
)

func DebugLog(c *gin.Context) {
	if conf.IsDebug() {
		c.Next()
		return
	}

	log.Debug("req:", c.Request.URL.String())
	log.Debug("body:", c.Request.Body)
}

func FillSession(c *gin.Context) {
	userId := header.UserId(c.Request)
	s, ok := session.GetUser(userId)
	if ok {
		c.Set(valuer.KeyUser, valuer.User{UserId: s.UserId, UserName: s.UserName})
	} else {
		//TODO
	}
}
