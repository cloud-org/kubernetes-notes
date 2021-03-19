package main

import (
	"kubernetes-notes/client/xterm_web_attach_container/server"
	"kubernetes-notes/client/xterm_web_attach_container/storage"
	"log"
	"net/http"
)

func main() {
	err := storage.InitKubeClient()
	if err != nil {
		log.Printf("init kube client err:%v\n", err)
		return
	}

	http.HandleFunc("/ssh", server.WsHandler)
	log.Printf("start server...")
	err = http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Printf("start server err:%v\n", err)
		return
	}
}
