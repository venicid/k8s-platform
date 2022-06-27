package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"k8s-platform/service"
	"net/http"
)

var Pod pod
type pod struct {

}


//Controller中的方法入参是gin.Context，用于从上下文中获取请求参数及定义响应内容
//流程：绑定参数->调用service代码->根据调用结果响应具体内容

// 获取Pod列表，支持分页，过滤，排序
func (p *pod) GetPods(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		FilterName string `form:"filter_name"`
		Namespace string `form:"namespace"`
		Limit int `form:"limit"`
		Page int `form:"page"`
	})

	// form格式使用bind方法，json格式使用showBuildJSON方法
	if err := ctx.Bind(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
		})
	
}


// 获取Pod详情
func (p *pod) GetPodDetail(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		PodName string `form:"pod_name"`
		Namespace string `form:"namespace"`
	
	})

	// form格式使用bind方法，json格式使用showBuildJSON方法
	if err := ctx.Bind(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod详情成功",
		"data": data,
	})

}



// 删除Pod
func (p *pod) DeletePod(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		PodName string `json:"pod_name"`
		Namespace string `json:"namespace"`

	})

	// json格式使用showBuildJSON方法
	if err := ctx.ShouldBindJSON(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Pod成功",
		"data": nil,
	})

}


// 更新Pod

func (p *pod) UpdatePod(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		PodName string `json:"pod_name"`
		Namespace string `json:"namespace"`
		Content string `json:"content"`  // 更新pod需要转义content

	})

	// json格式使用showBuildJSON方法
	if err := ctx.ShouldBindJSON(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	pod, err := service.Pod.UpdatePod(params.PodName, params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新Pod成功",
		"data": pod,
	})
}

/**
http://www.jsons.cn/jsonzip/ 转义工具
{
    "namespace":"default",
    "pod_name":"web-5d59997fc7-9jnsp",
    "content": "{\"metadata\":{\"name\":\"web-5d59997fc7-9jnsp\",\"generateName\":\"web-5d59997fc7-\",\"namespace\":\"default\",\"uid\":\"58099890-8a3e-4275-8128-b5b1a4825662\",\"resourceVersion\":\"1837276\",\"creationTimestamp\":\"2022-06-21T07:11:45Z\",\"labels\":{\"app\":\"web-pod\",\"pod-template-hash\":\"5d59997fc7\",\"test\":\"aaa\"},\"annotations\":{\"cni.projectcalico.org/containerID\":\"f6126366af6bf21053b5a89082cbac0fc3889693c9479e998f0d1c61624371ad\",\"cni.projectcalico.org/podIP\":\"10.182.127.247/32\",\"cni.projectcalico.org/podIPs\":\"10.182.127.247/32\"},\"ownerReferences\":[{\"apiVersion\":\"apps/v1\",\"kind\":\"ReplicaSet\",\"name\":\"web-5d59997fc7\",\"uid\":\"39509c60-419f-43a2-9741-84f6f64a3ae1\",\"controller\":true,\"blockOwnerDeletion\":true}],\"managedFields\":[{\"manager\":\"kube-controller-manager\",\"operation\":\"Update\",\"apiVersion\":\"v1\",\"time\":\"2022-06-21T07:11:45Z\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:metadata\":{\"f:generateName\":{},\"f:labels\":{\".\":{},\"f:app\":{},\"f:pod-template-hash\":{}},\"f:ownerReferences\":{\".\":{},\"k:{\\\"uid\\\":\\\"39509c60-419f-43a2-9741-84f6f64a3ae1\\\"}\":{\".\":{},\"f:apiVersion\":{},\"f:blockOwnerDeletion\":{},\"f:controller\":{},\"f:kind\":{},\"f:name\":{},\"f:uid\":{}}}},\"f:spec\":{\"f:containers\":{\"k:{\\\"name\\\":\\\"nginx\\\"}\":{\".\":{},\"f:image\":{},\"f:imagePullPolicy\":{},\"f:name\":{},\"f:resources\":{},\"f:terminationMessagePath\":{},\"f:terminationMessagePolicy\":{}}},\"f:dnsPolicy\":{},\"f:enableServiceLinks\":{},\"f:restartPolicy\":{},\"f:schedulerName\":{},\"f:securityContext\":{},\"f:terminationGracePeriodSeconds\":{}}}},{\"manager\":\"calico\",\"operation\":\"Update\",\"apiVersion\":\"v1\",\"time\":\"2022-06-21T07:11:48Z\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:metadata\":{\"f:annotations\":{\".\":{},\"f:cni.projectcalico.org/containerID\":{},\"f:cni.projectcalico.org/podIP\":{},\"f:cni.projectcalico.org/podIPs\":{}}}}},{\"manager\":\"kubelet\",\"operation\":\"Update\",\"apiVersion\":\"v1\",\"time\":\"2022-06-24T06:39:38Z\",\"fieldsType\":\"FieldsV1\",\"fieldsV1\":{\"f:status\":{\"f:conditions\":{\"k:{\\\"type\\\":\\\"ContainersReady\\\"}\":{\".\":{},\"f:lastProbeTime\":{},\"f:lastTransitionTime\":{},\"f:status\":{},\"f:type\":{}},\"k:{\\\"type\\\":\\\"Initialized\\\"}\":{\".\":{},\"f:lastProbeTime\":{},\"f:lastTransitionTime\":{},\"f:status\":{},\"f:type\":{}},\"k:{\\\"type\\\":\\\"Ready\\\"}\":{\".\":{},\"f:lastProbeTime\":{},\"f:lastTransitionTime\":{},\"f:status\":{},\"f:type\":{}}},\"f:containerStatuses\":{},\"f:hostIP\":{},\"f:phase\":{},\"f:podIP\":{},\"f:podIPs\":{\".\":{},\"k:{\\\"ip\\\":\\\"10.182.127.247\\\"}\":{\".\":{},\"f:ip\":{}}},\"f:startTime\":{}}}}]},\"spec\":{\"volumes\":[{\"name\":\"default-token-ltd25\",\"secret\":{\"secretName\":\"default-token-ltd25\",\"defaultMode\":420}}],\"containers\":[{\"name\":\"nginx\",\"image\":\"registry.cn-guangzhou.aliyuncs.com/chenxin_pub/nginx\",\"resources\":{},\"volumeMounts\":[{\"name\":\"default-token-ltd25\",\"readOnly\":true,\"mountPath\":\"/var/run/secrets/kubernetes.io/serviceaccount\"}],\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\",\"imagePullPolicy\":\"Always\"}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"serviceAccountName\":\"default\",\"serviceAccount\":\"default\",\"nodeName\":\"192.168.122.13\",\"securityContext\":{},\"schedulerName\":\"default-scheduler\",\"tolerations\":[{\"key\":\"node.kubernetes.io/not-ready\",\"operator\":\"Exists\",\"effect\":\"NoExecute\",\"tolerationSeconds\":300},{\"key\":\"node.kubernetes.io/unreachable\",\"operator\":\"Exists\",\"effect\":\"NoExecute\",\"tolerationSeconds\":300}],\"priority\":0,\"enableServiceLinks\":true,\"preemptionPolicy\":\"PreemptLowerPriority\"},\"status\":{\"phase\":\"Running\",\"conditions\":[{\"type\":\"Initialized\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"2022-06-21T07:11:45Z\"},{\"type\":\"Ready\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"2022-06-24T06:39:38Z\"},{\"type\":\"ContainersReady\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"2022-06-24T06:39:38Z\"},{\"type\":\"PodScheduled\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"2022-06-21T07:11:45Z\"}],\"hostIP\":\"192.168.122.13\",\"podIP\":\"10.182.127.247\",\"podIPs\":[{\"ip\":\"10.182.127.247\"}],\"startTime\":\"2022-06-21T07:11:45Z\",\"containerStatuses\":[{\"name\":\"nginx\",\"state\":{\"running\":{\"startedAt\":\"2022-06-24T06:39:36Z\"}},\"lastState\":{\"terminated\":{\"exitCode\":255,\"reason\":\"Error\",\"startedAt\":\"2022-06-21T07:11:50Z\",\"finishedAt\":\"2022-06-24T06:39:00Z\",\"containerID\":\"docker://7222e1ab4a1911ba03424983f3779ac91d8872fbb53a898312ae5866735202f3\"}},\"ready\":true,\"restartCount\":1,\"image\":\"registry.cn-guangzhou.aliyuncs.com/chenxin_pub/nginx:latest\",\"imageID\":\"docker-pullable://registry.cn-guangzhou.aliyuncs.com/chenxin_pub/nginx@sha256:3f13b4376446cf92b0cb9a5c46ba75d57c41f627c4edb8b635fa47386ea29e20\",\"containerID\":\"docker://b06b6caaf90908223a091b56e47dd6ab964a3f9561da9a11fc22da9749643f92\",\"started\":true}],\"qosClass\":\"BestEffort\"}}"
}
 */


// 获取Pod容器
func (p *pod) GetPodContainer(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		PodName string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})

	// json格式使用showBuildJSON方法
	if err := ctx.Bind(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	containers, err := service.Pod.GetPodContainer(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器成功",
		"data": containers,
	})
}


// 获取Pod中容器的日志
func (p *pod) GetPodContainerLog(ctx *gin.Context)  {
	// 处理入参
	params := new(struct{
		PodName string `form:"pod_name"`
		Namespace string `form:"namespace"`
		ContainerName string `form:"container_name"`
	})

	// json格式使用showBuildJSON方法
	if err := ctx.Bind(params); err != nil{
		logger.Error("bind绑定参数失败，"+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "bind绑定参数失败，"+err.Error(),
			"data": nil,
		})
		return
	}

	log, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器日志成功",
		"data": log,
	})
}


// 获取每个namepsace的pod数量
func (p *pod) GetPodNumPerNp(ctx *gin.Context)  {
	podsNps, err := service.Pod.GetPodNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取每个namepsace的pod数量",
		"data": podsNps,
	})
}