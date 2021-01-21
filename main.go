package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yumi/api"
	"yumi/apidoc"
	"yumi/conf"
	"yumi/gin"
	"yumi/gin/middleware"
	"yumi/pkg/log"
	"yumi/usecase"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	conf.Load()
	log.Init()

	log.Info("初始化数据库")
	//dbc.Init(conf.Get().DB)

	log.Info("初始化casbin")
	middleware.InitCasbin("", nil) //TODO:

	log.Info("初始化usecase")
	usecase.Init()

	log.Info("构建服务器")
	srvconf := conf.Get().Server
	mux := gin.NewMux()
	server := http.Server{
		Handler:      mux,
		Addr:         srvconf.Addr,
		ReadTimeout:  srvconf.ReadTimeout.Duration(),
		WriteTimeout: srvconf.WriteTimeout.Duration(),
	}
	mux.Use(
		middleware.Recovery(),
		middleware.Cors(conf.Get().CORS.AllowedOrigins, int(conf.Get().CORS.MaxAge.Duration()/time.Second)),
		middleware.Debug(conf.IsDebug()),
	)

	log.Info("加载路由")
	router := mux.Group("/")
	api.Mount(router)

	//debug模式下，开启接口文档
	if conf.IsDebug() {
		apidoc.Mount(router)
	}

	//启动服务
	log.Info("开始启动服务器，侦听地址：" + conf.Get().Server.Addr)
	go func() {
		if err := gin.Run(&server); err != nil {
			log.Info(fmt.Errorf("启动服务器失败: %s", err.Error()))
			c <- syscall.SIGINT
		}
	}()

end:
	for {
		s := <-c
		switch s {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			err := server.Shutdown(context.Background())
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
