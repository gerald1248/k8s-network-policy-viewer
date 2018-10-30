package main

import (
	//"encoding/json"
	"errors"
	"fmt"
	//"github.com/ghodss/yaml"
	"io/ioutil"
)

func processBytes(bytes []byte) (Result, error) {

	//preflight with optional conversion from YAMLs
	err := preflightAsset(&bytes)
	if err != nil {
		return Result{}, errors.New(fmt.Sprintf("input failed preflight check: %v", err))
	}

	//make sure config objects are presented as a list
	err = makeList(&bytes)
	if err != nil {
		return Result{}, err
	}

	result := Result{string(bytes[:])}

	return result, nil
}

func processFile(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("can't read %s: %v", path, err))
	}

	result, err := processBytes(bytes)

	if err != nil {
		return "", errors.New(fmt.Sprintf("can't process %s: %s", path, err))
	}

	return result.Buffer, nil
}
