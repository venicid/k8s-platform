package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Workflow workflow

type workflow struct {

}


// 获取d列表，支持分页，过滤，排序
func (w *workflow) GetWorkflows(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		FilterName string `form:"filter_name"`
		Namespace string `form:"namespace"`
		Limit int `form:"limit"`
		Page int `form:"page"`
	})

	// form格式使用bind方法，json格式使用showBuildJSON方法
	if err := ctx.Bind(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Workflow.GetWorkflows(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取workflows列表成功",
		"data": data,
	})

}


// 获取Pod详情
func (w *workflow) GetById(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		Id int `form:"id"`
	})

	// form格式使用bind方法，json格式使用showBuildJSON方法
	if err := ctx.Bind(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Workflow.GetById(params.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取workflow详情成功",
		"data": data,
	})

}



// 删除Pod
func (w *workflow) DeleteById(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		Id int `form:"id"`
	})

	// json格式使用showBuildJSON方法
	if err := ctx.Bind(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Workflow.DeleteById(params.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除workflow成功",
		"data": nil,
	})

}


// 新增workflow

func (w *workflow) CreateWorkFlow(ctx *gin.Context)  {

	var (
		wc = &service.WorkflowCreate{}
		err error
	)

	// json格式使用showBuildJSON方法
	if err := ctx.ShouldBindJSON(wc); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	 err = service.Workflow.CreateWorkFlow(wc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "新建Workflow成功",
		"data": nil,
	})
}

/** 创建数据
{
        "name": "tstng-wf064",
        "namespace": "default",
        "replicas": 1,
        "image":"nginx",
        "resource":"0.5/1",
        "health_check":false,
        "health_path":"",
        "label":{
            "app":"tstng-wf064"
        },
        "container_port":80,
        "type": "Ingress",
        "port":80,
        "node_port":null,
        "hots":"www.qwer.com",
        "path":"/",
        "path_type":"Prefix",
        "cpu":"0.5",
        "memory":"1Gi",
        "hosts":{
            "www.zzz.com":[
                {
                    "path":"/",
                    "path_type":"Prefix",
                    "service_name":"tstng-wf064",
                    "service_port":80
                }
            ]
        }
    }
*/