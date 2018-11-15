package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func processBytes(byteArray []byte, output *string) (string, int, error) {

	//preflight with optional conversion from YAMLs
	err := preflightAsset(&byteArray)
	if err != nil {
		return "", 0, errors.New(fmt.Sprintf("input failed preflight check: %v", err))
	}

	//make sure config objects are presented as a list
	err = makeList(&byteArray)
	if err != nil {
		return "", 0, err
	}

	var apiObjectSet ApiObjectSet

	if err = json.Unmarshal(byteArray, &apiObjectSet); err != nil {
		return "", 0, errors.New(fmt.Sprintf("can't unmarshal data: %v", err))
	}

	namespacePodMap := make(map[string][]string)
	podLabelMap := make(map[string]map[string]string)
	networkPolicyNamespaces := make(map[string]struct{})
	networkPolicies := []ApiObject{}
	for _, apiObject := range apiObjectSet.ApiObjects {
		// TODO: white/blacklist mechanism
		namespace := apiObject.Metadata.Namespace
		if strings.HasPrefix(namespace, "kube-") {
			continue
		}
		switch apiObject.Kind {
		case "Pod":
			// TODO: omitted skip condition: apiObject.Status.ContainerStatuses[0].Ready == true
			if len(apiObject.Status.ContainerStatuses) > 0 &&
				apiObject.Status.ContainerStatuses[0].Ready == true {
				namespacePodMap[namespace] = append(namespacePodMap[namespace], apiObject.Metadata.Name)
				podLabelMap[apiObject.Metadata.Name] = apiObject.Metadata.Labels
			}
		case "NetworkPolicy":
			networkPolicies = append(networkPolicies, *apiObject)
			networkPolicyNamespaces[namespace] = struct{}{}
		}
	}

	globalNamespaces := []string{}
	for podNamespace, _ := range namespacePodMap {
		if _, ok := networkPolicyNamespaces[podNamespace]; !ok {
			globalNamespaces = append(globalNamespaces, podNamespace)
		}
	}

	edgeMap := make(map[string][]string)
	initializeEdgeMap(&edgeMap, &namespacePodMap)
	allEdgesCount := countEdges(&edgeMap)
	filterEdgeMap(&edgeMap, &namespacePodMap, &podLabelMap, &networkPolicies, &globalNamespaces)
	filteredEdgesCount := countEdges(&edgeMap)

	var buffer bytes.Buffer
	switch *output {
	case "dot":
		writeDot(&namespacePodMap, &edgeMap, &buffer)
	case "json":
		writeJson(&namespacePodMap, &buffer)
	case "yaml":
		writeYaml(&namespacePodMap, &buffer)
	}
	var percentage float64
	percentage = (float64(filteredEdgesCount)/float64(allEdgesCount)) * 100.0
	percentageInt := int(percentage + 0.5)
	return buffer.String(), percentageInt, nil
}

func processFile(path string, output *string) (string, error) {
	byteArray, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("can't read %s: %v", path, err))
	}

	result, _, err := processBytes(byteArray, output)

	if err != nil {
		return "", errors.New(fmt.Sprintf("can't process %s: %s", path, err))
	}

	return result, nil
}

func countEdges(edgeMap *map[string][]string) int {
	count := 0
	for _, v := range (*edgeMap) {
		for _, _ = range v {
			count++
		}
	}
	return count
}
