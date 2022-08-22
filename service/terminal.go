package service

import (
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
)

var Terminal terminal

type terminal struct{}



// 消息内容
//TerminalMessage定义了终端和容器shell交互内容的格式
//Operation是操作类型
//Data是具体数据内容
//Rows和Cols可以理解为终端的行数和列数，也就是宽、高
type terminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}



//定义websocket的handler方法, 主流程
func (t *terminal) WhHandler(w http.ResponseWriter, r *http.Request)  {
	//加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil{
		logger.Error("创建k8s配置失败, " + err.Error())
		return
	}
	//解析form入参，获取namespace、podName、containerName参数
	if err := r.ParseForm(); err != nil{
		logger.Error("解析参数失败, " + err.Error())
		return
	}

	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	logger.Info("exec pod: %s, container: %s, namespace: %s\n", podName, containerName, namespace)

	// new一个TerminalSession类型的pty实例
	pty, err := NewTerminalSession(w , r, nil)
	if err != nil{
		logger.Error("实例化NewTerminalSession failed:%v\n", err)
		return
	}
	// 处理关闭
	defer func() {
		logger.Info("close session.")
		pty.Close()
	}()

	// 初始化pod所在的corev1资源组
	//PodExecOptions struct 包括Container stdout stdout Command 等结构
	//scheme.ParameterCodec 应该是pod 的GVK （GroupVersion & Kind）之类的
	//URL: // https://192.168.1.11:6443/api/v1/namespaces/default/pods/nginx-wf2-778d88d7c- 7rmsk/exec?command=%2Fbin%2Fbash&container=nginx- wf2&stderr=true&stdin=true&stdout=true&tty=true
	req := K8s.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").Name(podName).Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
		Container: containerName,
		Command:   []string{"bin/sh"},
	}, scheme.ParameterCodec)

	logger.Info("exec post request url:" , req.URL())

	// 升级SPDY协议
	//remotecommand 主要实现了http 转 SPDY 添加X-Stream-Protocol-Version相关header 并发送请求
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil{
		logger.Error("升级SPDY协议失败，" +  err.Error())
		return
	}
	// 与kubelet简历stream连接
	// 建立链接之后从请求的sream中发送、读取数据
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})
	if err != nil{
		logger.Error("执行pod命令失败，" +  err.Error())
		pty.Write([]byte("执行pod命令失败，" +  err.Error()))
		//标记退出stream流
		pty.Done()
	}
	
}


