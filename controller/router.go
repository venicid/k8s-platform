package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 初始化router类型的对象，首字母大写，用于跨包调用
var Router router

// 声明一个router的结构体
type router struct {}


// 初始化路由规则
func (r *router) InitApiRouter(router *gin.Engine){
	router.
		// 测试
		GET("/testapi", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "test success",
				"data": nil,
				})
			}).
		GET("/", func(c *gin.Context) {fmt.Printf("ClientIP: %s\n", c.ClientIP())}).

		// 登录
		POST("/api/login", Login.Auth).

		// 工作流
		POST("/api/k8s/workflow/create", Workflow.CreateWorkFlow).
		GET("/api/k8s/workflows", Workflow.GetWorkflows).
		GET("/api/k8s/workflow/detail",  Workflow.GetById).
		DELETE("/api/k8s/workflow/del", Workflow.DeleteById).

		// Pod操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodContainerLog).
		GET("/api/k8s/pod/numnp", Pod.GetPodNumPerNp).

		//namespace操作
		GET("/api/k8s/namespaces", Namespace.GetNamespaces).
		GET("/api/k8s/namespace/detail", Namespace.GetNamespaceDetail).
		DELETE("/api/k8s/namespace/del", Namespace.DeleteNamespace).

		// deployment
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		GET("/api/k8s/deployment/numnp", Deployment.GetDeployNumPerNp).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).

		//node操作
		GET("/api/k8s/nodes", Node.GetNodes).
		GET("/api/k8s/node/detail", Node.GetNodeDetail)

		//configmap操作
		//GET("/api/k8s/configmaps", ConfigMap.GetConfigMaps).
		//GET("/api/k8s/configmap/detail", ConfigMap.GetConfigMapDetail).
		//DELETE("/api/k8s/configmap/del", ConfigMap.DeleteConfigMap).
		//PUT("/api/k8s/configmap/update", ConfigMap.UpdateConfigMap)

		////sercret操作
		//GET("/api/k8s/secrets", Secret.GetSecrets).
		//GET("/api/k8s/secret/detail", Secret.GetSecretDetail).
		//DELETE("/api/k8s/secret/del", Secret.DeleteSecret).
		//PUT("/api/k8s/secret/update", Secret.UpdateSecret).
		//

		////daemonset操作
		//GET("/api/k8s/daemonsets", DaemonSet.GetDaemonSets).
		//GET("/api/k8s/daemonset/detail", DaemonSet.GetDaemonSetDetail).
		//DELETE("/api/k8s/daemonset/del", DaemonSet.DeleteDaemonSet).
		//PUT("/api/k8s/daemonset/update", DaemonSet.UpdateDaemonSet).
		//
		////statefulset操作
		//GET("/api/k8s/statefulsets", StatefulSet.GetStatefulSets).
		//GET("/api/k8s/statefulset/detail", StatefulSet.GetStatefulSetDetail).
		//DELETE("/api/k8s/statefulset/del", StatefulSet.DeleteStatefulSet).
		//PUT("/api/k8s/statefulset/update", StatefulSet.UpdateStatefulSet).
		//
		////service操作
		//GET("/api/k8s/services", Servicev1.GetServices).
		//GET("/api/k8s/service/detail", Servicev1.GetServiceDetail).
		//DELETE("/api/k8s/service/del", Servicev1.DeleteService).
		//PUT("/api/k8s/service/update", Servicev1.UpdateService).
		//POST("/api/k8s/service/create", Servicev1.CreateService).
		//
		////ingress操作
		//GET("/api/k8s/ingresses", Ingress.GetIngresses).
		//GET("/api/k8s/ingress/detail", Ingress.GetIngressDetail).
		//DELETE("/api/k8s/ingress/del", Ingress.DeleteIngress).
		//PUT("/api/k8s/ingress/update", Ingress.UpdateIngress).
		//POST("/api/k8s/ingress/create", Ingress.CreateIngress).

		////pvc操作
		//GET("/api/k8s/pvcs", Pvc.GetPvcs).
		//GET("/api/k8s/pvc/detail", Pvc.GetPvcDetail).
		//DELETE("/api/k8s/pvc/del", Pvc.DeletePvc).
		//PUT("/api/k8s/pvc/update", Pvc.UpdatePvc).
		//
		//
		////pv操作
		//GET("/api/k8s/pvs", Pv.GetPvs).
		//GET("/api/k8s/pv/detail", Pv.GetPvDetail)

}