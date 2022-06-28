package main

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/db"
	"k8s-platform/service"
)

func main(){

	// 初始化数据库
	db.Init()

	// 初始化k8s client
	service.K8s.Init()   // 可以使用server.k8s.ClientSet

	// 初始化gin
	r := gin.Default()
	// 跨包调用router的初始化方法
	controller.Router.InitApiRouter(r)

	// 启动gin server
	r.Run(config.ListenAddr)

	// 关闭数据库连接
	db.Close()

}




