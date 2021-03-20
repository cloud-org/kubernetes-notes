package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"kubernetes-notes/client/util"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
)

var podName string
var container string

func init() {
	flag.StringVar(&podName, "pods", "", "pod name")
	flag.StringVar(&container, "c", "", "container name")
	flag.Usage = usage
}

func main() {
	flag.Parse()
	if podName == "" || container == "" {
		//fmt.Println("podName or container should not be nil")
		flag.Usage()
		return
	}
	clientset, err := util.InitClient()
	if err != nil {
		return
	}
	//var tailLines int64 = 100

	namespace := "default"

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{
		Container: container,
		Follow:    true,
		//TailLines: &tailLines,
	})
	fmt.Println(req.URL())

	ctx := context.TODO()

	readCloser, err := req.Stream(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer readCloser.Close()

	r := bufio.NewReader(readCloser)
	go func() {
		for {
			bytes, err := r.ReadBytes('\n')
			fmt.Println(string(bytes)) // TODO: websocket
			if err != nil {
				//if err != io.EOF {
				//	fmt.Println("")
				//	return
				//}
				fmt.Println("read err", err.Error())
				return
			}
		}
	}()

	time.Sleep(5 * time.Second)
	return

}

func usage() {
	fmt.Fprintf(os.Stdout, `pod_logs - get pods log
Usage: pod_logs [-h help]
Options:
`)
	flag.PrintDefaults()
}
