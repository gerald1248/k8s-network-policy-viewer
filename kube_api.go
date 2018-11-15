package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"strings"
)

// don't allow errors to bubble up
// TODO: switch to byte array as that's what is used in the end
func getJsonData(buffer *string) {
	*buffer = "{}"
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Can't create in-cluster configuration: %s\n", err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Can't create clientset: %s\n", err.Error())
		return
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Pods not found\n")
		return
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing pods: %v\n", statusError.ErrStatus.Message)
		return
	} else if err != nil {
		log.Printf("API error: %s\n", err.Error())
		return
	}

	// TODO: is there an option to retain the "Kind" field on retrieval?
	for index, _ := range pods.Items {
		pods.Items[index].Kind = "Pod"
	}

	podsJson, err := json.Marshal(pods.Items)
	if err != nil {
		log.Printf("Can't marshal pods: %s\n", err.Error())
		return
	}

	networkPolicies, err := clientset.NetworkingV1().NetworkPolicies("").List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		log.Printf("Network policies not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Printf("Error listing network policies: %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		log.Printf("API error: %s\n", err.Error())
		return
	}

	for index, _ := range networkPolicies.Items {
		networkPolicies.Items[index].Kind = "NetworkPolicy"
	}

	networkPoliciesJson, err := json.Marshal(networkPolicies.Items)
	if err != nil {
		log.Printf("Can't marshal network policies: %s\n", err.Error())
		return
	}

	podsJsonString := string(podsJson)
	networkPoliciesJsonString := string(networkPoliciesJson)
	trimBrackets(&podsJsonString)
	trimBrackets(&networkPoliciesJsonString)

	// TODO: switch to go templating

	*buffer = fmt.Sprintf("{\"kind\":\"List\",\"apiVersion\":\"v1\",\"Items\":[%s,%s]}", podsJsonString, networkPoliciesJsonString)
}

func trimBrackets(s *string) {
	*s = strings.Trim(*s, "[]")
}
