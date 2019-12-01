package tail

import (
	"easyexec/pkg/common"
	"easyexec/pkg/login"
	"easyexec/pkg/objs"
	"log"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type Tail struct {
	common.CommonType
	FilePath string
}

func NewTail(f string) *Tail {
	return &Tail{
		FilePath: f,
	}
}

func (t *Tail) Tail() {
	var line = "20"

	// var podName string = "kube-apiserver-docker-for-desktop"
	kubeconfig, err := login.GetConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	podGetter := objs.NewPodGetter(clientset)
	podGetter.PodName = t.PodName
	podGetter.DeployName = t.DeploymentName
	podGetter.ContainerIndex = t.ContainerIndex
	var pod *v1.Pod
	if t.Namespace != "" {
		podGetter.Namespace = t.Namespace
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

	ExecuteRemoteCommand(clientset, clientConfig, pod, "tail -f -n "+line+" "+t.FilePath)
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
