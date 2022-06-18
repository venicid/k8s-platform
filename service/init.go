package service

import (

	"k8s-platform/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct{
	 client *kubernetes.Clientset
}

func (k *k8s) InitClient()  {
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil{
		panic(err)
	}

	// 根据rest.config类型的对象，new 一个clientset出来
	var clientset  *kubernetes.Clientset
	clientset, err = kubernetes.NewForConfig(conf)
	if err != nil{
		panic(err)
	}

	k.client = clientset

}