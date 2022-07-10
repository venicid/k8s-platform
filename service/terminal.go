package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"time"
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

// ws交互结构体，接管输入和输出
//定义TerminalSession结构体，实现PtyHandler接口
//wsConn是websocket连接
//sizeChan用来定义终端输入和输出的宽和高
//doneChan用于标记退出终端
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
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
		Command:   []string{"bin/bash"},
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


//初始化一个websocket.Upgrader类型的对象，用于http协议升级为websocket协议
//var upgrader = &websocket.Upgrader{
//	HandshakeTimeout:  time.Second * 2,
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return upgrader
}()

//该方法用于升级http协议至websocket，并new一个TerminalSession类型的对象返回
func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {

	// 升级ws
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}

	return session, nil
}

//用于读取web端的输入，接收web端输入的指令内容
func (t *TerminalSession) Read(p []byte) (int, error) {

	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		log.Printf("read message err: %v", err)
		return 0, err
	}

	var msg terminalMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		log.Printf("read parse message err: %v", err)
		return 0, nil
	}

	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{
			Width:  msg.Cols,
			Height: msg.Rows,
		}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("unknown message type '%s'", msg.Operation)
		return 0, fmt.Errorf("unknown message type '%s'", msg.Operation)
	}

}

//用于向web端输出，接收web端的指令后，将结果返回出去
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("write parse message err: %v", err)
		return 0, err
	}

	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("write  message err: %v", err)
		return 0, err
	}

	return len(p), nil
}

//获取web端是否resize，以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}

}

// 关闭doneChan，关闭后触发退出终端
func (t *TerminalSession) Done() {
	close(t.doneChan)

}

// 用于关闭websocket连接
func (t *TerminalSession) Close() error {
	return t.wsConn.Close()

}


