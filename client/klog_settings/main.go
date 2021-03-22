package main

import (
	"flag"
	"fmt"

	"k8s.io/klog/v2"
)

func main() {

	klog.InitFlags(nil)

	err := flag.Set("stderrthreshold", "4")
	if err != nil {
		fmt.Println(err)
	}
	flag.Parse()

	klog.Infof("hello %s", "world")

	klog.V(5).Info("需要 -v 5 或者以上才可以打印出来!")

	defer klog.Flush()

	return
}
