package tailf

import (
	"easyexec/pkg/common"
	"easyexec/pkg/objs"
	"encoding/json"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
	"strconv"
)

type Tail struct {
	common.CommonType
	FilePath string
	Count    int
}

func NewTail(f string, count int) *Tail {
	return &Tail{
		FilePath: f,
		Count:    count,
	}
}

func (t *Tail) Tail() {
	clientConfig, clientset := common.GetKubeClient()
	podGetter := objs.NewPodGetter(clientset)
	podGetter.PodName = t.PodName
	podGetter.DeployName = t.DeploymentName
	podGetter.ContainerIndex = t.ContainerIndex
	var pod *v1.Pod
	var err error
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

	t.executeRemoteCommand(clientset, clientConfig, pod, "tail -f -n "+strconv.Itoa(t.Count)+" "+t.FilePath)
}

func (t *Tail) executeRemoteCommand(coreClient kubernetes.Interface, restCfg *restclient.Config, pod *v1.Pod, command string) {
	if t.execute(coreClient, restCfg, pod, "sh", command) != nil {
		err := t.execute(coreClient, restCfg, pod, "bash", command)
		if err != nil {
			fmt.Println(err)
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
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	j, err := json.Marshal(err)
	fmt.Println(2, string(j))
	return err
}
