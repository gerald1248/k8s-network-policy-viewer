package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func processBytes(byteArray []byte, output *string) (Result, error) {

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
	switch *output {
	case "dot":
		writeDot(&namespacePodMap, &buffer)
	}
	return Result{buffer.String()}, nil
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

	return result.Buffer, nil
}
