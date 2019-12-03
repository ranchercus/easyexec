package list

import (
	"easyexec/pkg/common"
	"easyexec/pkg/objs"
	"k8s.io/api/core/v1"
	"log"
	"strings"
)

type podListResult struct {
	Name      string
	Namespace string
}

type podList struct {
	common.CommonType
}

func NewPodList() *podList {
	return &podList{}
}

func (p *podList) List() []podListResult {
	_, clientset := common.GetKubeClient()
	podgetter := objs.NewPodGetter(clientset)
	var nss []string
	var err error
	if p.Namespace != "" {
		nss = []string{p.Namespace}
	} else {
		nss, err = objs.GetNsByCookie()
		if err != nil {
			log.Fatal(err)
		}
	}
	podgetter.PodName = p.PodName
	podgetter.DeployName = p.DeploymentName

	result := make([]podListResult, 0)
	for _, ns := range nss {
		podgetter.Namespace = ns
		list, err := podgetter.List()
		if err != nil {
			continue
		}
		r := p.compose(list)
		result = append(result, r...)
	}
	return result
}

func (p *podList) compose(list []*v1.Pod) []podListResult {
	result := make([]podListResult, 0)
	for _, v := range list {
		r := &podListResult{
			v.Name,
			v.Namespace,
		}
		result = append(result, *r)
	}
	return result
}

func (p *podList) List2String() string {
	result := p.List()
	s := make([]string, 0)
	s = append(s, "命名空间\t\t名称")
	for _, p := range result {
		s = append(s, p.Namespace+"\t\t"+p.Name)
	}
	return strings.Join(s, "\r\n")
}
