package util

import (
	"io/ioutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// TODO: 可配置
const kubeConfigPath = "/home/ashing/.kube/kind-config-third-k8s-cluster"

// 初始化 k8s 客户端
func InitClient() (clientset *kubernetes.Clientset, err error) {
	var (
		restConf *rest.Config
	)

	if restConf, err = GetRestConf(); err != nil {
		return
	}

	// 生成 clientset 配置
	if clientset, err = kubernetes.NewForConfig(restConf); err != nil {
		goto END
	}
END:
	return
}

// 获取 k8s restful client 配置
func GetRestConf() (restConf *rest.Config, err error) {
	var (
		kubeconfig []byte
	)

	// 读 kubeconfig 文件
	if kubeconfig, err = ioutil.ReadFile(kubeConfigPath); err != nil {
		goto END
	}
	// 生成 rest client 配置
	if restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		goto END
	}
END:
	return
}
