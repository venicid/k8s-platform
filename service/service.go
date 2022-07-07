package service

import (
	"context"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//定义ServiceCreate结构体，用于创建service需要的参数属性的定义
type ServiceCreate struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
	Type string `json:"type"`
	ContainerPort int32 `json:"container_port"`
	Port int32 `json:"port"`
	NodePort int32 `json:"node_port"`
	Label map[string]string `json:"label"`
}

var ServiceV1 servicev1

type servicev1 struct {}

//创建service,,接收ServiceCreate对象
func (s *servicev1) CreateService(data *ServiceCreate) (err error){
	//将data中的数据组装成corev1.Service对象
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       data.Name,
			Namespace:                  data.Namespace,
			Labels:                     data.Label,
			},

		//Spec中定义类型，端口，选择器
		Spec:       corev1.ServiceSpec{
			Type:                          corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name: "https",
					Port: data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
		Status:     corev1.ServiceStatus{},
	}

	//默认ClusterIP,这里是判断NodePort,添加配置
	if data.NodePort != 0 && data.Type =="NodePort"{
		service.Spec.Ports[0].NodePort = data.NodePort
	}

	// 创建service
	_,err = K8s.ClientSet.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil{
		logger.Error("创建ingress失败，"+err.Error())
		return  errors.New("创建ingress失败，"+err.Error())
	}

	return nil

}

/**
apiVersion: v1
kind: Service
metadata:
	name: myapp-svc
	namespace: default
spec:
	selector:
		app: myapp
	ports:
	- name: http
		port: 80 #service的端口
		protocol： tcp #协议
		targetPort: 80 #pod的端口
 */

// 删除service
func (s *servicev1) DeleteService(serviceName, namespace string) ( err error)  {
	err = K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if  err!= nil{
		logger.Error("删除service失败," + err.Error())
		return  errors.New("删除service失败," + err.Error())
	}
	return nil
}