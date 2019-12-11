package common

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func GetKubeClient() (*restclient.Config, *kubernetes.Clientset) {
	kubeconfig, err := GetConfigPath()
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
	return clientConfig, clientset
}

func GetSAKubeClient() (*restclient.Config, *kubernetes.Clientset) {
	clientConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error){
		return clientcmd.Load([]byte(KUBE_SA_CONFIG))
	})
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Fatal(err)
	}
	return clientConfig, clientset
}