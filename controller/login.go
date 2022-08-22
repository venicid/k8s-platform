package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Login login

type login struct {

}

//验证账号密码
func (l *login) Auth(ctx *gin.Context)  {
	params := new(struct{
					Username string `json:"username"`
					Password string `json:"password"`
					})
	if err := ctx.ShouldBindJSON(params); err != nil{
		logger.Error("bind绑定参数失败" +  err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败" +  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Login.Auth(params.Username, params.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
		"data": nil,
	})

}
