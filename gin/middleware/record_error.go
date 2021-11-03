package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"

	"yumi/gin"
	"yumi/pkg/log"
	"yumi/pkg/types"
)

// RecordError ...
func RecordError() gin.HandlerFunc {
	var warpErr types.WarpError = "middleware.RecordError"

	return func(c *gin.Context) {
		c.Next()
		if c.Error != nil {
			// 所有接口的错误日志都在这里打印
			bindBuf := &bytes.Buffer{}
			err := json.NewEncoder(bindBuf).Encode(c.BindObject)
			if err != nil {
				log.Warn(warpErr.Warp(err))
			}
			content := []interface{}{
				fmt.Sprintln("fullpath:", c.FullPath()),
				fmt.Sprintln("request uri:", c.Request.URL.Path),
				fmt.Sprintln("request params:", c.Params),
				fmt.Sprintln("request bindobject:", bindBuf.String()),
				fmt.Sprintln("error:", c.Error),
			}
			log.Error(content...)
		}
	}
}
