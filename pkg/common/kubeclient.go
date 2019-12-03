package common

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
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
