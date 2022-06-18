package service

import (
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)


// 用于初始化k8s client
var K8s k8s

type k8s struct{
	 ClientSet *kubernetes.Clientset
}

// 初始化方法
func (k *k8s) Init()  {
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil{
		panic("获取k8s client配置失败" + err.Error())
	}

	// 根据rest.config类型的对象，new 一个clientset出来
	var clientSet  *kubernetes.Clientset
	clientSet, err = kubernetes.NewForConfig(conf)
	if err != nil{
		panic("创建k8s client失败" + err.Error())
	}else{
		logger.Info("k8s client 初始化成功！")
	}

	k.ClientSet = clientSet

}