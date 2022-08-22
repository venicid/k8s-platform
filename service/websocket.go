package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"net/http"
	"time"
)

const END_OF_TRANSMISSION = "\u0004"

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

// Stdin ...
func (t *TerminalSession) Stdin() io.Reader {
	return t
}

// Stdout ...
func (t *TerminalSession) Stdout() io.Writer {
	return t
}

// Stderr ...
func (t *TerminalSession) Stderr() io.Writer {
	return t
}

