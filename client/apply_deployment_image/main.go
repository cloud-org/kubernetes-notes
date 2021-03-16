package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kubernetes-notes/client/util"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	//"k8s.io/api/apps/v1beta1"
	appv1 "k8s.io/api/apps/v1"
)

func main() {

	var err error
	clientset, err := util.InitClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 读取 yaml
	deploymentYaml, err := ioutil.ReadFile("./config/nginx_deployment.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// yaml 转 json
	deploymentJson, err := yaml.ToJSON(deploymentYaml)
	if err != nil {
		fmt.Println(err)
		return
	}

	var deployment appv1.Deployment
	// json 转 struct
	err = json.Unmarshal(deploymentJson, &deployment)
	if err != nil {
		fmt.Println(err)
		return
	}

	nginxContainer := v1.Container{}
	nginxContainer.Name = "nginx"
	nginxContainer.Image = "nginx:1.14.1"
	containers := make([]v1.Container, 0)
	containers = append(containers, nginxContainer)

	deployment.Spec.Template.Spec.Containers = containers

	// 查询 k8s 是否有该 deployment
	namespace := "default"
	ctx := context.TODO()
	_, err = clientset.AppsV1().Deployments(namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		// 如果是除了未找到的其他错误 直接退出
		if !errors.IsNotFound(err) {
			return
		}
		// 不存在则创建
		_, err = clientset.AppsV1().Deployments(namespace).Create(ctx, &deployment, metav1.CreateOptions{})
		if err != nil {
			return
		}
	} else { // 存在则更新
		_, err = clientset.AppsV1().Deployments(namespace).Update(ctx, &deployment, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("apply deployment 成功")
	return

}
