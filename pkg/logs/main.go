package main

import (
	"bufio"
	"fmt"
	"io"
	"show_logs/objs"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

var line = int64(20)

func openStream(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) (io.ReadCloser, error) {
	return client.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(podID).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).Stream()
}

func GetLogFile(client kubernetes.Interface, namespace, podID string, container string, usePreviousLogs bool) (io.ReadCloser, error) {
	logOptions := &v1.PodLogOptions{
		Container:  container,
		Follow:     true,
		Previous:   usePreviousLogs,
		Timestamps: false,
		TailLines:  &line,
	}
	logStream, err := openStream(client, namespace, podID, logOptions)
	return logStream, err
}

func main() {
	var podName string = "kube-apiserver-docker-for-desktop"
	kubeconfig := "/Users/Tibbers/.kube/config"
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	pod := objs.GetAvaliablePod(clientset, podName)
	logStream, err := GetLogFile(clientset, pod.Namespace, pod.Name, pod.Spec.Containers[0].Name, false)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(logStream)
	for {
		line, _, err := r.ReadLine()
		if err == nil {
			fmt.Println(string(line))
		}
	}
}
