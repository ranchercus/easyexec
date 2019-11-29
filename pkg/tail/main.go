package main

import (
	"easyexec/pkg/objs"
	"fmt"
	"log"
	"os"
	"show_logs/objs"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func showHelp() {
	fmt.Println("")
	fmt.Println("使用方法：")
	fmt.Println("logexec 容器ID [可选参数...]")
	fmt.Println("")
	fmt.Println("可选参数：")
	fmt.Println("	日志路径， 默认/home/tomcat/目录下的app.log日志")
	fmt.Println("	显示最后多少行日志，默认最后20行")
	fmt.Println("如：logexec containerid /home/tomcat/log/*/app.log 30")
	fmt.Println("")
}

func main() {
	var podName string
	var path = "/home/tomcat/log/*/app.log"
	var line = "20"
	args := os.Args
	if len(args) < 2 {
		showHelp()
		return
	} else if len(args) == 2 {
		podName = args[1]
	} else if len(args) == 3 {
		podName = args[1]
		path = args[2]
	} else if len(args) == 4 {
		podName = args[1]
		path = args[2]
		line = args[3]
	} else {
		showHelp()
		return
	}

	// var podName string = "kube-apiserver-docker-for-desktop"
	kubeconfig := "/Users/Tibbers/.kube/config"
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	pod := objs.GetAvaliablePod(clientset, podName)
	ExecuteRemoteCommand(clientset, clientConfig, pod, "tail -f -n "+line+" "+path)
}

func ExecuteRemoteCommand(coreClient kubernetes.Interface, restCfg *restclient.Config, pod *v1.Pod, command string) {
	request := coreClient.CoreV1().
		RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: strings.Split(command, " "),
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		log.Fatal(err)
	}
}
