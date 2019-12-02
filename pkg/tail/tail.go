package tail

import (
	"easyexec/pkg/common"
	"easyexec/pkg/login"
	"easyexec/pkg/objs"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
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

	t.executeRemoteCommand(clientset, clientConfig, pod, "tail -f -n "+line+" "+t.FilePath)
}

func (t *Tail) executeRemoteCommand(coreClient kubernetes.Interface, restCfg *restclient.Config, pod *v1.Pod, command string) {
	if t.execute(coreClient, restCfg, pod, "sh", command) != nil {
		err := t.execute(coreClient, restCfg, pod, "bash", command)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (t *Tail) execute(coreClient kubernetes.Interface, restCfg *restclient.Config, pod *v1.Pod, scmd, command string) error {
	request := coreClient.CoreV1().
		RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: pod.Spec.Containers[t.ContainerIndex-1].Name,
			Command:   []string{scmd, "-c", command},
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
			Stdin:     true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	})
	return err
}
