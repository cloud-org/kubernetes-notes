package storage

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//WsMessage websocket 消息
type WsMessage struct {
	MessageType int
	Data        []byte
}

type Ws interface {
	ReadLoop()
	WriteLoop()
	Write(msgType int, data []byte) (err error)
	Read() (msg *WsMessage, err error)
	Close()
}

//WsConnection websocket 连接
//closChan 关闭之后 还是可以获取到零值 https://colobu.com/2016/04/14/Golang-Channels/#close
type WsConnection struct {
	wsSocket *websocket.Conn
	inChan   chan *WsMessage // receive msg
	outChan  chan *WsMessage // send msg

	mu        sync.Mutex
	isClosed  bool
	closeChan chan byte
}

//NewWsConnection
//Upgrade 升级为 websocket 连接
func NewWsConnection(resp http.ResponseWriter, req *http.Request) (wsConn *WsConnection, err error) {
	// 允许所有跨域请求
	wsUpgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	sock, err := wsUpgrader.Upgrade(resp, req, nil) //important
	if err != nil {
		return
	}
	wsConn = &WsConnection{
		wsSocket:  sock,
		inChan:    make(chan *WsMessage, 1000),
		outChan:   make(chan *WsMessage, 1000),
		mu:        sync.Mutex{},
		isClosed:  false,
		closeChan: make(chan byte),
	}

	// 启动协程
	go wsConn.ReadLoop()
	go wsConn.WriteLoop()

	return
}

//ReadLoop 循环读取前端发送的消息并写入 inChan
func (w *WsConnection) ReadLoop() {
	for {
		msgType, data, err := w.wsSocket.ReadMessage()
		if err != nil {
			log.Println(err)
			w.Close()
			return
		}
		msg := &WsMessage{
			MessageType: msgType,
			Data:        data,
		}
		select {
		case w.inChan <- msg:
		case <-w.closeChan:
			log.Println("read loop close chan")
			return
		}
	}
}

//WriteLoop 循环读取 outChan 中的消息并发送到前端
func (w *WsConnection) WriteLoop() {
	for {
		select {
		case msg := <-w.outChan:
			err := w.wsSocket.WriteMessage(msg.MessageType, msg.Data)
			if err != nil { // TODO: 是否要忽略错误
				log.Println(err)
				return
			}
		case <-w.closeChan:
			log.Println("write loop close chan")
			return
		}
	}
}

//Write
func (w *WsConnection) Write(msgType int, data []byte) (err error) {
	select {
	case w.outChan <- &WsMessage{msgType, data}:
	case <-w.closeChan:
		err = errors.New("websocket closed when write")
	}
	return
}

//Read
func (w *WsConnection) Read() (msg *WsMessage, err error) {
	select {
	case msg = <-w.inChan:
		return
	case <-w.closeChan:
		err = errors.New("websocket closed when read")
	}
	return
}

//Close
func (w *WsConnection) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.isClosed {
		err := w.wsSocket.Close()
		if err != nil {
			log.Printf("close ws conn err: %v\n", err)
		}
		w.isClosed = true
		close(w.closeChan)
	}
}
