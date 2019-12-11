package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var template = `package common

var KUBE_SA_CONFIG = %s%s%s`

var config string

func main() {
	flag.StringVar(&config, "kubeconfig", "", "kubernetes cluster config file path")
	flag.Parse()

	if config == "" {
		log.Fatal("kubeconfig cant be null")
	}

	bfile, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(err)
	}

	kubeconfig := fmt.Sprintf(template, "`", string(bfile), "`")

	err = ioutil.WriteFile("pkg/common/saconfig.go", []byte(kubeconfig), 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generate config OK.")
}
