package logs

import (
	"bufio"
	"easyexec/pkg/common"
	"easyexec/pkg/objs"
	"fmt"
	"io"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
)

type Logs struct {
	common.CommonType
}

func NewLogs(commonType common.CommonType) *Logs {
	return &Logs{
		commonType,
	}
}

var line = int64(20)

func (l *Logs) openStream(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) (io.ReadCloser, error) {
	return client.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(podID).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).Stream()
}

func (l *Logs) Logs() {
	_, clientset := common.GetKubeClient()
	podGetter := objs.NewPodGetter(clientset)
	podGetter.PodName = l.PodName
	podGetter.DeployName = l.DeploymentName
	podGetter.ContainerIndex = l.ContainerIndex
	var pod *v1.Pod
	var err error
	if l.Namespace != "" {
		podGetter.Namespace = l.Namespace
		pod, err = podGetter.Get()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		nss, err := objs.GetNsByCookie()
		if err != nil {
			log.Fatal(err)
		}
		for _, ns := range nss {
			podGetter.Namespace = ns
			pod, err = podGetter.Get()
			if err == nil {
				break
			}
		}
		if pod == nil {
			log.Fatal("无法智能查询到POD")
		}
	}

	logOptions := &v1.PodLogOptions{
		Container:  pod.Spec.Containers[l.ContainerIndex-1].Name,
		Follow:     true,
		Previous:   false,
		Timestamps: false,
		TailLines:  &line,
	}
	logStream, err := l.openStream(clientset, pod.Namespace, pod.Name, logOptions)
	if err != nil {
		log.Fatal("无法连接上POD", err)
	}
	r := bufio.NewReader(logStream)
	for {
		line, _, err := r.ReadLine()
		if err == nil {
			fmt.Println(string(line))
		}
	}
}
