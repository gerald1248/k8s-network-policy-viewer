package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"fmt"
)

func processBytes(byteArray []byte) (Result, error) {

	//preflight with optional conversion from YAMLs
	err := preflightAsset(&byteArray)
	if err != nil {
		return Result{}, errors.New(fmt.Sprintf("input failed preflight check: %v", err))
	}

	//make sure config objects are presented as a list
	err = makeList(&byteArray)
	if err != nil {
		return Result{}, err
	}

	var apiObjectSet ApiObjectSet

	if err = json.Unmarshal(byteArray, &apiObjectSet); err != nil {
		return Result{}, errors.New(fmt.Sprintf("can't unmarshal data: %v", err))
	}

	namespacePodMap := map[string][]string{}
	for _, apiObject := range apiObjectSet.ApiObjects {
		if apiObject.Kind != "Pod" {
			continue
		}
		namespacePodMap[apiObject.Metadata.Namespace] = append(namespacePodMap[apiObject.Metadata.Namespace], apiObject.Metadata.Name)
	}

	var buffer bytes.Buffer
	buffer.WriteString("digraph podNetwork {\n")
	for k, v := range namespacePodMap {
		buffer.WriteString("  subgraph cluster_")
		buffer.WriteString(k)
		buffer.WriteString(" {\n")
		for _, s := range v {
			buffer.WriteString("    \"")
			buffer.WriteString(s)
			buffer.WriteString("\";\n")
		}
		buffer.WriteString("  }\n")
	}
	buffer.WriteString("}\n")
	return Result{buffer.String()}, nil
}

func processFile(path string) (string, error) {
        byteArray, err := ioutil.ReadFile(path)
        if err != nil {
                return "", errors.New(fmt.Sprintf("can't read %s: %v", path, err))
        }

        result, err := processBytes(byteArray)

        if err != nil {
                return "", errors.New(fmt.Sprintf("can't process %s: %s", path, err))
        }

        return result.Buffer, nil
}
