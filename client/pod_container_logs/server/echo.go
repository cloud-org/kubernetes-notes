package server

import (
	"context"
	"fmt"
	"kubernetes-notes/client/pod_container_logs/storage"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "k8s.io/api/core/v1"
)

//WsHandler ssh controller
func echoWsHandler(c echo.Context) error {
	podNs := c.QueryParam("podNs")
	podName := c.QueryParam("podName")
	containerName := c.QueryParam("containerName")

	//podNs := "default"
	//podName := "nginx-deployment-66b6c48dd5-9bw5x"
	//containerName := "nginx"

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

	// websocket readCloser handler struct
	returnChan := make(chan bool)
	logStreamHandler := storage.NewLogStreamHandler(readCloser, wsConn, returnChan)
	logStreamHandler.Handle()

	<-returnChan
	log.Println("ws api return")

	return nil

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
