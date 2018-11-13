package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

func getJsonData(buffer *string) {
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
	podsJson, err := json.Marshal(pods)
	if err != nil {
		log.Printf("Can' marshal pods: %s", err.Error())
	}

	/*
	networkPolicies, err := clientset.CoreV1().NetworkPolicies("").List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Network policies not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing network policies: %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	}
	networkPoliciesJson, err := json.Marshal(networkPolicies)
	if err != nil {
		log.Printf("Can' marshal network policies: %s", err.Error())
	}
	*/

	*buffer = fmt.Sprintf("%v", podsJson)
}
