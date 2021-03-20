package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"kubernetes-notes/client/pod_container_logs/storage"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "k8s.io/api/core/v1"
)

//XtermMessage
//resize 字段类型参考 remotecommand.TerminalSize
type XtermMessage struct {
	MsgType string `json:"msgtype"`
	Input   string `json:"input"` // use when msg type is input
	Rows    uint16 `json:"rows"`  // use when msg type is resize
	Cols    uint16 `json:"cols"`  // use when msg type is resize
}

//WsHandler ssh controller
func echoWsHandler(c echo.Context) error {
	//podNs := c.QueryParam("podNs")
	//podName := c.QueryParam("podName")
	//containerName := c.QueryParam("containerName")

	podNs := "default"
	podName := "nginx-deployment-66b6c48dd5-9bw5x"
	containerName := "nginx"

	wsConn, err := storage.NewWsConnection(c.Response(), c.Request())
	if err != nil {
		log.Printf("init ws connect err:%v\n", err)
		return c.JSON(http.StatusInternalServerError, "error")
	}
	// 之后如果发生错误，需要将 wsConn 连接关闭
	defer wsConn.Close()

	req := storage.KubeClient.CoreV1().Pods(podNs).GetLogs(podName, &v1.PodLogOptions{
		Container: containerName,
		Follow:    true,
	})
	fmt.Println(req.URL())

	readCloser, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer readCloser.Close()

	// TODO: websocket readCloser handler struct

	r := bufio.NewReader(readCloser)
	go func() {
		for {
			bytes, err := r.ReadBytes('\n')
			if err != nil {
				log.Println("read err", err.Error())
				return
			}
			//fmt.Println(string(bytes))
			//dst := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
			//base64.StdEncoding.Encode(dst, bytes)
			err = wsConn.Write(websocket.TextMessage, bytes)
			if err != nil {
				log.Println("ws write msg err", err.Error())
				return
			}
		}
	}()

	msg, err := wsConn.Read()
	if err != nil {
		log.Println("ws read err", err)
		return nil
	}
	var xtermMessage XtermMessage
	if err = json.Unmarshal(msg.Data, &xtermMessage); err != nil {
		log.Println(err)
		return nil
	}

	return nil
	//	var
	//	for {
	//		msg, err = wsConn.Read()
	//		if err != nil {
	//			return
	//		}
	//		if err = json.Unmarshal(msg.Data, &xtermMessage); err != nil {
	//			return
	//		}
	//		log.Printf("xterm msg is %+v\n", xtermMessage)
	//		switch xtermMessage.MsgType {
	//		case INPUT:
	//			size = len(xtermMessage.Input)
	//			copy(p, xtermMessage.Input)
	//		case RESIZE:
	//			sh.resizeEvent <- remotecommand.TerminalSize{
	//				Width:  xtermMessage.Cols,
	//				Height: xtermMessage.Rows,
	//			}
	//		}
	//	}
	//}()

	//<- closeChan // TODO: if goroutine1 or goroutine2 exist

	//return nil
}

func addMiddleware(e *echo.Echo) {
	// 增加 cors 中间件
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}

func addApi(e *echo.Echo) {
	e.GET("/ssh", echoWsHandler)
}

//CreateEngine echo
func CreateEngine() (*echo.Echo, error) {
	e := echo.New()

	addMiddleware(e)
	addApi(e)

	return e, nil
}
