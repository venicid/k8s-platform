package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Node node
type node struct {}

//获取node列表，支持过滤、排序、分页
func (n *node) GetNodes(ctx *gin.Context)  {
	// 参数
	params := new(struct{
		FilterName string `form:"filter_name"`
		Page int `form:"page"`
		Limit int `form:"limit"`
	})
	if err := ctx.Bind(params); err!=nil{
		logger.Error("Bind请求参数失败,", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
	}
	// 请求
	data, err := service.Node.GetNodes(params.FilterName, params.Limit, params.Page)
	if err !=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
	}

	// 返回
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Node列表成功",
		"data": data,
	})
}


//获取node详情
func (n *node) GetNodeDetail(ctx *gin.Context)  {
	// 参数
	params := new (struct{
		NodeName string `form:"node_name"`
	})
	if err := ctx.Bind(params); err!= nil{
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	// 请求service
	data, err  := service.Node.GetNodeDetail(params.NodeName)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	// 返回
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Node详情成功",
		"data": data,
	})

}