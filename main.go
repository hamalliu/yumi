package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"yumi/api"
	"yumi/apidoc"
	"yumi/pkg/conf"
	"yumi/pkg/log"
	"yumi/pkg/gin"
	"yumi/pkg/gin/middeware"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	conf.Load()
	log.Init()

	log.Info("初始化数据库")
	//dbc.Init(conf.Get().DB)

	log.Info("构建服务器")
	srv := gin.NewServer()
	srv.Use(middeware.Recovery(), middeware.Cors(conf.Get().CORS), middeware.Debug())

	log.Info("加载路由")
	router := srv.Group("/")
	api.Mount(router)

	//debug模式下，开启接口文档
	if conf.IsDebug() {
		apidoc.Mount(router)
	}

	//启动服务
	log.Info("开始启动服务器，侦听地址：" + conf.Get().Server.Addr)
	go func() {
		if err := srv.Run(conf.Get().Server); err != nil {
			log.Info(fmt.Errorf("启动服务器失败: %s", err.Error()))
			c <- syscall.SIGINT
		}
	}()

end:
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			err := srv.Shutdown(context.Background())
			if err != nil {
				log.Info("关闭服务器失败:", err.Error())
			} else {
				log.Info("服务器已关闭")
			}
			break end
		case syscall.SIGHUP:
			break
		default:
			break
		}
	}
}
