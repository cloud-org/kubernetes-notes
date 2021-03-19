package storage

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

type Handler interface {
	Next() *remotecommand.TerminalSize
	Read(p []byte) (size int, err error)
	Write(p []byte) (size int, err error)
}

//StreamHandler ssh 流式处理器
type StreamHandler struct {
	wsConn      *WsConnection
	resizeEvent chan remotecommand.TerminalSize
}

func NewStreamHandler(wsConn *WsConnection, resizeEvent chan remotecommand.TerminalSize) *StreamHandler {
	return &StreamHandler{wsConn: wsConn, resizeEvent: resizeEvent}
}

//XtermMessage
//resize 字段类型参考 remotecommand.TerminalSize
type XtermMessage struct {
	MsgType string `json:"msgtype"`
	Input   string `json:"input"` // use when msg type is input
	Rows    uint16 `json:"rows"`  // use when msg type is resize
	Cols    uint16 `json:"cols"`  // use when msg type is resize
}

func (sh *StreamHandler) Next() *remotecommand.TerminalSize {
	res := <-sh.resizeEvent
	return &res
}

//Read 读取前端输入
func (sh *StreamHandler) Read(p []byte) (size int, err error) {
	var (
		msg          *WsMessage
		xtermMessage XtermMessage
	)
	msg, err = sh.wsConn.Read()
	if err != nil {
		return
	}
	if err = json.Unmarshal(msg.Data, &xtermMessage); err != nil {
		return
	}
	log.Printf("xterm msg is %+v\n", xtermMessage)
	switch xtermMessage.MsgType {
	case INPUT:
		size = len(xtermMessage.Input)
		copy(p, xtermMessage.Input)
	case RESIZE:
		sh.resizeEvent <- remotecommand.TerminalSize{
			Width:  xtermMessage.Cols,
			Height: xtermMessage.Rows,
		}
	}
	return
}

//Write 向 web 输出
func (sh *StreamHandler) Write(p []byte) (size int, err error) {
	log.Printf("write p is %v\n", string(p))
	size = len(p) //size 就是 len(p)
	b := []byte(base64.StdEncoding.EncodeToString(p))
	copyData := make([]byte, len(b))
	copy(copyData, b)
	err = sh.wsConn.Write(websocket.TextMessage, copyData)
	return
}
