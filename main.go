package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/dao"
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

	/**
	测试
	 */
	// 测试workflow数据库连接
	data, _ := dao.Workflow.GetLWorkflows("nginx", "default", 10, 1)
	fmt.Println(data)

	res, _ := dao.Workflow.GetById(3)
	fmt.Println(res)


	// 启动gin server
	r.Run(config.ListenAddr)

	// 关闭数据库连接
	db.Close()

}




