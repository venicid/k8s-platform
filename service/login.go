package service

import (
	"errors"
	"github.com/wonderivan/logger"
	"k8s-platform/config"
)

var Login login

type login struct {

}

func (l *login) Auth(userName, password string)  (err error){
	if userName == config.AdminUser&& password == config.AdminPwd {
		return nil
	}
	logger.Error("登录失败，用户名or密码错误")
	return errors.New("登录失败，用户名or密码错误")
}