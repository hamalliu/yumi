package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"yumi/api"
	"yumi/api_doc"
	"yumi/pkg/conf"
	"yumi/pkg/log"
	"yumi/pkg/net/ymhttp"
	"yumi/pkg/net/ymhttp/middeware"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	conf.Load()
	log.Init()

	log.Info("初始化数据库")
	//dbc.Init(conf.Get().DB)

	log.Info("构建服务")
	srv := ymhttp.NewServer()
	srv.Use(middeware.Recovery(), middeware.Cors(conf.Get().CORS), middeware.PrintRequest())

	log.Info("加载路由")
	router := srv.Group("/")
	api.Mount(router)

	//debug模式下，开启接口文档
	if conf.IsDebug() {
		api_doc.Mount(router)
	}

	//启动服务
	log.Info("开始启动服务，侦听地址：" + conf.Get().Server.Addr)
	go func() {
		if err := srv.Run(conf.Get().Server); err != nil {
			log.Info(fmt.Errorf("启动服务失败: %s", err.Error()))
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
				log.Info("server shutdown error:", err.Error())
			} else {
				log.Info("server shutdown success")
			}
			break end
		case syscall.SIGHUP:
			break
		default:
			break
		}
	}
}
