package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 初始化router类型的对象，首字母大写，用于跨包调用
var Router router

// 声明一个router的结构体
type router struct {}


// 初始化路由规则
func (r *router) InitApiRouter(router *gin.Engine){
	router.GET("/testapi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "test success",
			"data": nil,
		})
	}).
	GET("/", func(c *gin.Context) {fmt.Printf("ClientIP: %s\n", c.ClientIP())}).
		POST("/api/k8s/workflow/create", Workflow.CreateWorkFlow).
		GET("/api/k8s/workflows", Workflow.GetWorkflows).
		GET("/api/k8s/workflow/detail",  Workflow.GetById).
		DELETE("/api/k8s/workflow/del", Workflow.DeleteById).
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodContainerLog).
		GET("/api/k8s/pod/numnp", Pod.GetPodNumPerNp)




}