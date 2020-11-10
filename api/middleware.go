package api

import (
	"yumi/apistorage/session"
	"yumi/pkg/conf"
	"yumi/pkg/log"
	"yumi/pkg/net/gin"
	"yumi/pkg/net/gin/header"
	"yumi/pkg/net/gin/valuer"
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

//FillSession ...
func FillSession(c *gin.Context) {
	userID := header.UserID(c.Request)
	s, ok := session.GetUser(userID)
	if ok {
		c.Set(valuer.KeyUser, valuer.User{UserID: s.UserID, UserName: s.UserName})
	} else {
		//TODO
	}
}
