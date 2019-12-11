package exec

import (
	"easyexec/pkg/common"
	"easyexec/pkg/objs"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"os"
)

type Exec struct {
	common.CommonType
	Cmd string
}

func NewExec(commonType common.CommonType, cmd string) *Exec {
	return &Exec{
		CommonType: commonType,
		Cmd:        cmd,
	}
}

func (t *Exec) Exec() {
	_, clientset := common.GetKubeClient()
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

	t.executeRemoteCommand(pod)
}

func (t *Exec) executeRemoteCommand(pod *v1.Pod) {
	if t.execute(pod, "sh") != nil {
		err := t.execute(pod, "bash")
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
}

func (t *Exec) execute(pod *v1.Pod, scmd string) error {
	restCfg, coreClient := common.GetSAKubeClient()
	request := coreClient.CoreV1().
		RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: pod.Spec.Containers[t.ContainerIndex-1].Name,
			Command:   []string{scmd, "-c", t.Cmd},
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
	return err
}
