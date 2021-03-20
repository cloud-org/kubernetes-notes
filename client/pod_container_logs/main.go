package main

import (
	"context"
	"flag"
	"fmt"
	"kubernetes-notes/client/util"
	"os"

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
	res := req.Do(ctx)
	if res.Error() != nil {
		fmt.Println(res.Error())
		return
	}
	b, err := res.Raw()
	if err != nil {
		fmt.Println(b)
		return
	}
	fmt.Println(string(b))
}

func usage() {
	fmt.Fprintf(os.Stdout, `pod_logs - get pods log
Usage: pod_logs [-h help]
Options:
`)
	flag.PrintDefaults()
}
