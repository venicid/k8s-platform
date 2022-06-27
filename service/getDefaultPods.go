package service

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getDefaultPods(){

	// 将kubeconfig文件转为res.config类型的对象
	// config 在 ls ~/.kube/config
	//conf, err := clientcmd.BuildConfigFromFlags("", "E:\\goProject\\config")
	// 公有k8s仓库
	conf, err := clientcmd.BuildConfigFromFlags("", "E:\\goProject\\kubectl.kubeconfig")
	if err != nil{
		fmt.Println("err")
		panic(err)
	}

	// 根据rest.config类型的对象，new 一个clientset出来
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil{
		panic(err)
	}

	//使用clientset获取pod列表
	podList, err := clientset.CoreV1().Pods("default").
		List(context.TODO(), metav1.ListOptions{})
	if err != nil{
		panic(err)
	}

	for _,pod := range podList.Items{
		fmt.Println(pod.Name, pod.Namespace)
	}

}

