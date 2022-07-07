package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct {

}

//定义ServiceCreate结构体，用于创建service需要的参数属性的定义
type IngressCreate struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Label map[string]string `json:"label"`
	Hosts map[string][]*HttpPath `json:"hosts"`
}

//定义ingress的path结构体
type HttpPath struct {
	Path string `json:"path"`
	PathType nwv1.PathType `json:"path_type"`
	ServiceName string `json:"service_name"`
	ServicePort int32 `json:"service_port"`
}

// 创建ingress
func (i *ingress) CreateIngress(data *IngressCreate)  (err error) {
	//声明nwv1.IngressRule和nwv1.HTTPIngressPath变量，后面组装数据
	var ingressRules []nwv1.IngressRule
	var httpIngressPATHs []nwv1.HTTPIngressPath

	//将data中的数据组装成nwv1.Ingress对象
	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       data.Name,
			Namespace:                  data.Namespace,
			Labels:                     data.Label,
		},
		Status:     nwv1.IngressStatus{},
	}

	// 第一层for循环是将host组装成nwv1.IngressRule类型的对象
	// 一个host对应一个ingressrule，每个ingressrule中包含一个host和多个path
	for key,value :=range data.Hosts{
		ir := nwv1.IngressRule{
			Host:             key,
			//这里现将nwv1.HTTPIngressRuleValue类型中的Paths置为空，后面组装好数据再赋值
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths: nil},
					},
		}

		//第二层for循环是将path组装成nwv1.HTTPIngressPath类型的对象
		for _,httpPath := range value{
			hip := nwv1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend:  nwv1.IngressBackend{
					Service:  &nwv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}

			//将每个hip对象组装成数组
			httpIngressPATHs = append(httpIngressPATHs, hip)
		}

		//给Paths赋值，前面置为空了
		ir.IngressRuleValue.HTTP.Paths = httpIngressPATHs
		//将每个ir对象组装成数组，这个ir对象就是IngressRule，每个元素是一个host和多个path
		ingressRules = append(ingressRules, ir)
	}

	//将ingressRules对象加入到ingress的规则中
	ingress.Spec.Rules = ingressRules

	// 创建ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil{
		logger.Error("创建ingress失败，"+err.Error())
		return  errors.New("创建ingress失败，"+err.Error())
	}


	return nil


}


/**
对标yaml

apiVersion: extensions/v1beta1
kind: Ingress
metadata:
	name: ingress-myapp
	namespace: default
spec:
	rules:
		- host: www.xxx.com
		  http:
			paths:
			- path: /
			  pathType: Prefix
			  backend:
				service:
					name: myapp-svc
					port:
						number: 80
 */


// 删除ingress
func (i *ingress) DeleteIngress(ingressName, namespace string) ( err error)  {
	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if  err!= nil{
		logger.Error("删除ingress失败," + err.Error())
		return  errors.New("删除ingress失败," + err.Error())
	}
	return nil
}