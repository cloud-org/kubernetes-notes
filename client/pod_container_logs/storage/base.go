package storage

import (
	"kubernetes-notes/client/util"

	"k8s.io/client-go/kubernetes"
)

const (
	RESIZE = "resize"
	INPUT  = "input"
)

var (
	KubeClient *kubernetes.Clientset
)

//InitKubeClient
func InitKubeClient() error {
	client, err := util.InitClient()
	if err != nil {
		return err
	}
	KubeClient = client
	return nil
}
