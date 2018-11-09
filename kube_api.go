package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

func listPods(buffer *string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})

	if errors.IsNotFound(err) {
		log.Printf("Pods not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing pods: %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	}

	*buffer = fmt.Sprintf("%v", pods)
}
