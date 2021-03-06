package server

import (
	"kubernetes-notes/client/util"
	"kubernetes-notes/client/xterm_web_attach_container/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

//WsHandler ssh controller
func ginWsHandler(c *gin.Context) {
	podNs := c.Query("podNs")
	podName := c.Query("podName")
	containerName := c.Query("containerName")

	wsConn, err := storage.NewWsConnection(c.Writer, c.Request)
	if err != nil {
		log.Printf("init ws connect err:%v\n", err)
		c.JSON(http.StatusInternalServerError, "init websocket error")
		return
	}
	// 之后如果发生错误，需要将 wsConn 连接关闭
	defer wsConn.Close()

	restConf, err := util.GetRestConf()
	if err != nil {
		log.Printf("get rest conf err:%v\n", err)
		return
	}

	sshReq := storage.KubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		SubResource("exec").
		Name(podName).
		Namespace(podNs).
		VersionedParams(&v1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: containerName,
			Command:   []string{"env", "LANG=C.UTF-8", "/bin/bash"}, // 支持中文
		}, scheme.ParameterCodec)

	log.Printf("ssh url is %s\n", sshReq.URL().String())

	executor, err := remotecommand.NewSPDYExecutor(restConf, "POST", sshReq.URL())
	if err != nil {
		log.Printf("init executor err:%v\n", err)
		return
	}

	handler := storage.NewStreamHandler(wsConn, make(chan remotecommand.TerminalSize))
	// 注: Stream 这里会阻塞
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		Tty:               true,
		TerminalSizeQueue: handler,
	})
	if err != nil {
		log.Printf("fix stream handler err:%v\n", err)
		return
	}

	//终端 exit 之后会走到这里
	return
}

//addGinApi
func addGinApi(g *gin.Engine) {
	g.GET("/ssh", ginWsHandler)
}

//CreateGinEngine
func CreateGinEngine() (*gin.Engine, error) {
	g := gin.Default()

	addGinApi(g)

	return g, nil
}
