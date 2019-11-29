package objs

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetAvaliablePod(clientset kubernetes.Interface, podName string) *v1.Pod {
	podClient := clientset.CoreV1().Pods("")
	pods, err := podClient.List(metav1.ListOptions{
		FieldSelector: "metadata.name=" + podName,
	})
	if err != nil {
		panic(err)
	}
	if len(pods.Items) == 0 {
		panic("无法找到POD")
	}
	var pod *v1.Pod
	for _, p := range pods.Items {
		if len(p.Spec.Containers) == 0 {
			continue
		}
		status := true
		for _, cs := range p.Status.ContainerStatuses {
			if !cs.Ready {
				status = false
				break
			}
		}
		if !status {
			continue
		}
		pod = &p
		break
	}
	if pod == nil {
		panic("POD状态异常")
	}
	return pod
}
