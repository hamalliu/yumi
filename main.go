package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"yumi/api"
	"yumi/conf"
	"yumi/gin"
	"yumi/gin/middleware"
	"yumi/pkg/log"
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/pkg/stores/mgoc"
	"yumi/usecase"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	conf.Load()
	log.InitStdlog()

	log.Info("构建mongodb客户端")
	mgoConf := conf.Get().Mongo
	mgoCli, err := mgoc.New(mgoConf.Dsn, mgoConf.Options()...)
	if err != nil {
		panic(err)
	}

	log.Info("构建mysqldb客户端")
	dbConf := conf.Get().DB
	myCLi, err := mysqlx.New(dbConf.Dsn, dbConf.Options()...)
	if err != nil {
		panic(err)
	}

	log.Info("初始化casbin")
	middleware.InitCasbin("", nil) //TODO:

	log.Info("初始化usecase")
	usecase.Init(mgoCli, myCLi)

	log.Info("构建服务器")
	srvconf := conf.Get().HttpServer
	mux := gin.NewMux()
	server := http.Server{
		Handler:      mux,
		Addr:         srvconf.Addr,
		ReadTimeout:  srvconf.ReadTimeout.Duration(),
		WriteTimeout: srvconf.WriteTimeout.Duration(),
	}

	log.Info("加载路由")
	router := mux.Group("yumi", "/")
	api.Mount(router)

	//启动服务
	log.Info("启动服务器，侦听地址：" + conf.Get().HttpServer.Addr)
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
		default:
		}
	}
}
