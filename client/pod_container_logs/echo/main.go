package main

import (
	"kubernetes-notes/client/pod_container_logs/server"
	"kubernetes-notes/client/pod_container_logs/storage"
	"log"
)

func main() {

	err := storage.InitKubeClient()
	if err != nil {
		log.Printf("init kube client err:%v\n", err)
		return
	}

	e, err := server.CreateEngine()
	if err != nil {
		log.Printf("create engine err: %v\n", err)
		return
	}

	log.Printf("start server on port 7777...")
	err = e.Start(":7777")
	if err != nil {
		log.Printf("start server err:%v\n", err)
		return
	}
}
