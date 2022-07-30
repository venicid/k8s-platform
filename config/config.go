package config

import "time"

const(
	ListenAddr = "0.0.0.0:9090"
	//KubeConfig = "E:\\goProject\\config"
	KubeConfig = "E:\\goProject\\admin.conf"
	//KubeConfig = "E:\\goProject\\kubectl.kubeconfig"

	// tail的日志行数
	// tail -n 2000
	PodLogTailLine = 2000

	// 数据库配置
	DbType = "mysql"
	DbHost = "150.158.171.205"
	DbPort = 3306
	DbName = "k8s_demo"
	DbUser = "root"
	DbPwd = "root112358"

	// 打印mysql debug的sql日志
	LogMode = true

	// 连接池的配置
	MaxIdleConns = 10 // 最大空闲连接
	MaxOpenConns = 100 // 最大连接数
	MaxLifeTime = 30 * time.Second // 最大生存时间

	// 用户名密码
	AdminUser = "admin"
	AdminPwd = "admin"
)
