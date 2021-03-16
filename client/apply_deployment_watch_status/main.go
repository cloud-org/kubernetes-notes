package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kubernetes-notes/client/util"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"

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

	// 给 pod 添加 label
	//deployment.Spec.Template.Labels["deploy_time"] = time.Now().In(time.Local).Format(time.RFC3339)
	deployment.Spec.Template.Labels["deploy_time"] = strconv.Itoa(int(time.Now().Unix())) // 123123

	namespace := "default"

	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		k8sDeployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}

		if k8sDeployment.Status.UpdatedReplicas == *(k8sDeployment.Spec.Replicas) &&
			k8sDeployment.Status.Replicas == *(k8sDeployment.Spec.Replicas) &&
			k8sDeployment.Status.AvailableReplicas == *(k8sDeployment.Spec.Replicas) &&
			k8sDeployment.Status.ObservedGeneration == k8sDeployment.Generation {
			// 滚动升级完成
			fmt.Println("滚动升级完成")
			break
		}

		// 打印工作中的 pod 比例
		fmt.Printf("部署中：(%d/%d)\n", k8sDeployment.Status.AvailableReplicas, *(k8sDeployment.Spec.Replicas))

	}

	fmt.Println("apply deployment 成功")

	// 打印每个 pod 的状态(可能会打印出 terminating 中的 pod, 但最终只会展示新 pod 列表)
	podList, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "app=nginx"})
	if err == nil {
		for _, pod := range podList.Items {
			podName := pod.Name
			podStatus := string(pod.Status.Phase)

			// PodRunning means the pod has been bound to a node and all of the containers have been started.
			// At least one container is still running or is in the process of being restarted.
			if podStatus == string(corev1.PodRunning) {
				// 汇总错误原因不为空
				if pod.Status.Reason != "" {
					podStatus = pod.Status.Reason
					goto KO
				}

				// condition 有错误信息
				for _, cond := range pod.Status.Conditions {
					if cond.Type == corev1.PodReady { // POD 就绪状态
						if cond.Status != corev1.ConditionTrue { // 失败
							podStatus = cond.Reason
						}
						goto KO
					}
				}

				// 没有 ready condition, 状态未知
				podStatus = "Unknown"
			}

		KO:
			fmt.Printf("[name:%s status:%s]\n", podName, podStatus)
		}
	}

	return

}
