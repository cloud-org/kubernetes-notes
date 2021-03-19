package main

import (
	"kubernetes-notes/client/xterm_web_attach_container/server"
	"kubernetes-notes/client/xterm_web_attach_container/storage"
	"log"
)

func main() {

	err := storage.InitKubeClient()
	if err != nil {
		log.Printf("init kube client err:%v\n", err)
		return
	}

	g, err := server.CreateGinEngine()
	if err != nil {
		log.Printf("create engine err: %v\n", err)
		return
	}

	//log.Printf("start server on port 7777...")
	err = g.Run(":7777")
	if err != nil {
		log.Printf("start server err:%v\n", err)
		return
	}
}
