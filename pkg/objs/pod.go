package objs

import (
	"errors"
	"strings"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodGetter struct {
	clientset      kubernetes.Interface
	PodName        string
	DeployName     string
	Namespace      string
	ContainerIndex int
}

func NewPodGetter(clientset kubernetes.Interface) *PodGetter {
	return &PodGetter{
		clientset: clientset,
	}
}

func (g *PodGetter) Get() (*v1.Pod, error) {
	var pods []v1.Pod
	var err error
	if g.DeployName == "" {
		if g.PodName == "" {
			return nil, errors.New("POD名不能为空")
		}
		pods, err = g.getPodListByName()
	} else {
		pods, err = g.getPodListByDeployment()
	}
	if err != nil {
		return nil, err
	}
	return g.filterAvaliablePod(pods)
}

func (g *PodGetter) getPodListByName() ([]v1.Pod, error) {
	podClient := g.clientset.CoreV1().Pods(g.Namespace)
	pods, err := podClient.List(metav1.ListOptions{
		FieldSelector: "metadata.name=" + g.PodName,
	})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func (g *PodGetter) getPodListByDeployment() ([]v1.Pod, error) {
	deployClient := g.clientset.AppsV1().Deployments(g.Namespace)
	deploy, err := deployClient.Get(g.DeployName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	podClient := g.clientset.CoreV1().Pods(g.Namespace)
	pods := make([]v1.Pod, 0)
	if g.PodName == "" {
		lbs := make([]string, 0)
		for k, v := range deploy.Spec.Selector.MatchLabels {
			lbs = append(lbs, k+"="+v)
		}
		plist, err := podClient.List(metav1.ListOptions{
			LabelSelector: strings.Join(lbs, ","),
		})
		if err != nil {
			return nil, err
		}
		pods = append(pods, plist.Items...)
	} else {
		pod, err := podClient.Get(g.PodName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		notFound := false
		for dlk, dlv := range deploy.Spec.Selector.MatchLabels {
			found := false
			for plk, plv := range pod.Labels {
				if plk == dlk && plv == dlv {
					found = true
					break
				}
			}
			if !found {
				notFound = true
			}
		}
		if notFound {
			return nil, errors.New("无法在Deployment中获取该POD")
		}
		pods = append(pods, *pod)
	}
	return pods, nil
}

func (g *PodGetter) filterAvaliablePod(pods []v1.Pod) (*v1.Pod, error) {
	if len(pods) == 0 {
		return nil, errors.New("无法找到POD")
	}
	var pod *v1.Pod
	for _, p := range pods {
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
		return nil, errors.New("POD状态异常")
	}
	return pod, nil
}
