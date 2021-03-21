package storage

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"

	"github.com/gorilla/websocket"
)

//XtermMessage
//resize 字段类型参考 remotecommand.TerminalSize
type XtermMessage struct {
	MsgType string `json:"msgtype"`
	Input   string `json:"input"` // use when msg type is input
	Rows    uint16 `json:"rows"`  // use when msg type is resize
	Cols    uint16 `json:"cols"`  // use when msg type is resize
}

type LogStreamHandler struct {
	readCloser io.ReadCloser
	wsConn     *WsConnection
	returnChan chan bool
}

func NewLogStreamHandler(readCloser io.ReadCloser, wsConn *WsConnection, returnChan chan bool) *LogStreamHandler {
	return &LogStreamHandler{readCloser: readCloser, wsConn: wsConn, returnChan: returnChan}
}

//Read message from frontend
func (l *LogStreamHandler) Read() {
	for {
		msg, err := l.wsConn.Read()
		if err != nil {
			log.Println("read from frontend err:", err)
			select {
			case l.returnChan <- true:
			default:
				log.Println("returnChan full")
			}
			return
		}
		var xtermMessage XtermMessage
		if err = json.Unmarshal(msg.Data, &xtermMessage); err != nil {
			log.Println(err)
			return
		}
		log.Println(xtermMessage)
	}
}

//Write write msg to frontend
func (l *LogStreamHandler) Write() {
	reader := bufio.NewReader(l.readCloser)
	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println("reader read bytes err:", err.Error())
			select {
			case l.returnChan <- true:
			default:
				log.Println("returnChan full")
			}
			return
		}
		dst := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
		base64.StdEncoding.Encode(dst, bytes)
		err = l.wsConn.Write(websocket.TextMessage, dst)
		if err != nil {
			log.Println("ws write msg err:", err.Error())
			return
		}
	}
}

func (l *LogStreamHandler) Handle() {
	go l.Read()
	go l.Write()
}
