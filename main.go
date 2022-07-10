package main

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/db"
	"k8s-platform/middle"
	"k8s-platform/service"
	"net/http"
)

func main(){

	// 初始化数据库
	db.Init()

	// 初始化k8s client
	service.K8s.Init()   // 可以使用server.k8s.ClientSet

	// 初始化gin
	r := gin.Default()

	// 跨域配置
	r.Use(middle.Cors())

	// jwt.Token
	//r.Use(middle.JWTAuth())

	// 跨包调用router的初始化方法
	controller.Router.InitApiRouter(r)

	/**
	测试
	 */
	// 测试workflow数据库连接
	//data, _ := dao.Workflow.GetLWorkflows("nginx", "default", 10, 1)
	//fmt.Println(data)
	//
	//res, _ := dao.Workflow.GetById(3)
	//fmt.Println(res)


	/**
	测试终端websocket
	地址：https://tool.uvooc.com/websocket/
	ws://localhost:8081/ws?namespace=default&pod_name=tstng-wf061-b965d784b-2v68r&container_name=tstng-wf061

	结果：
	连接成功，现在你可以发送信息进行测试了！
	服务端回应 2022-07-10 19:12:53
	{"operation":"stdout","data":"\u001b[?2004hroot@tstng-wf061-b965d784b-2v68r:/# ","rows":0,"cols":0}
	你发送的信息 2022-07-10 19:13:15
	{"operation":"stdin","data":"ls","rows":0,"cols":0}
	服务端回应 2022-07-10 19:13:16
	{"operation":"stdout","data":"ls","rows":0,"cols":0}
	你发送的信息 2022-07-10 19:13:29
	{"operation":"stdin","data":"\r","rows":0,"cols":0}
	服务端回应 2022-07-10 19:13:29
	{"operation":"stdout","data":"\r\n\u001b[?2004l\r","rows":0,"cols":0}
	服务端回应 2022-07-10 19:13:29
	{"operation":"stdout","data":"bin docker-entrypoint.d home media proc\tsbin tmp\r\nboot docker-entrypoint.sh lib
	mnt\t root\tsrv usr\r\ndev etc\t\t lib64 opt\t run\tsys var\r\n\u001b[?2004hroot@tstng-wf061-b965d784b-2v68r:/# ",
	"rows":0,"cols":0}
	 */
	go func() {
		http.HandleFunc("/ws", service.Terminal.WhHandler)
		http.ListenAndServe(":8081", nil)
	}()


	// 启动gin server
	r.Run(config.ListenAddr)

	// 关闭数据库连接
	db.Close()

}




