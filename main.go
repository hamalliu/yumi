package main

import (
	"fmt"

	"yumi/api"
	"yumi/conf"
	"yumi/controller"
	"yumi/utils/log"
)

func main() {
	forover := make(chan int)
	log.Info("开始初始化服务")

	log.Info("初始化数据库")
	//db.Init(conf.GetDB())

	log.Info("初始化控制器")
	controller.Init()
	log.Info("加载路由")
	api.Mount(controller.Route{Pattern: "/"})

	//启动服务
	log.Info("开始启动服务，侦听地址：" + conf.Get().Addr)
	go func() {
		if err := controller.Run(); err != nil {
			log.Info(fmt.Errorf("启动服务失败: %s", err.Error()))
			forover <- 1
		}
	}()

	<-forover
}
