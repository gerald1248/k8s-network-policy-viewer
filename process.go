package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func processBytes(byteArray []byte, output *string) (string, error) {

	//preflight with optional conversion from YAMLs
	err := preflightAsset(&byteArray)
	if err != nil {
		return "", errors.New(fmt.Sprintf("input failed preflight check: %v", err))
	}

	//make sure config objects are presented as a list
	err = makeList(&byteArray)
	if err != nil {
		return "", err
	}

	var apiObjectSet ApiObjectSet

	if err = json.Unmarshal(byteArray, &apiObjectSet); err != nil {
		return "", errors.New(fmt.Sprintf("can't unmarshal data: %v", err))
	}

	namespacePodMap := make(map[string][]string)
	podLabelMap := make(map[string]map[string]string)
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
		}
	}

	edgeMap := make(map[string][]string)
	initializeEdgeMap(&edgeMap, &namespacePodMap)
	filterEdgeMap(&edgeMap, &namespacePodMap, &podLabelMap, &networkPolicies)

	var buffer bytes.Buffer
	switch *output {
	case "dot":
		writeDot(&namespacePodMap, &edgeMap, &buffer)
	case "json":
		writeJson(&namespacePodMap, &buffer)
	case "yaml":
		writeYaml(&namespacePodMap, &buffer)
	}
	return buffer.String(), nil
}

func processFile(path string, output *string) (string, error) {
	byteArray, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("can't read %s: %v", path, err))
	}

	result, err := processBytes(byteArray, output)

	if err != nil {
		return "", errors.New(fmt.Sprintf("can't process %s: %s", path, err))
	}

	return result, nil
}
