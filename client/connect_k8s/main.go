package main

import (
	"context"
	"fmt"
	"kubernetes-notes/client/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	var (
		clientset *kubernetes.Clientset
		podsList  *corev1.PodList
		singlePod corev1.Pod
		index     int
		err       error
	)

	// 初始化 k8s 客户端
	if clientset, err = util.InitClient(); err != nil {
		goto FAIL
	}

	// 获取 default 命名空间下的所有 POD
	if podsList, err = clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{}); err != nil {
		goto FAIL
	}
	//fmt.Println(*podsList)
	for index, singlePod = range podsList.Items {
		fmt.Printf("index is %d, pod name is %v, pod phase is %v, pod ip is %v, host ip is %v\n",
			index, singlePod.Name, singlePod.Status.Phase, singlePod.Status.PodIP, singlePod.Status.HostIP)
	}

	return

FAIL:
	fmt.Println(err)
	return
}
