package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"unicode/utf8"
)

// ensure YAML as well as JSON can be read
// applies only to file-based processing; the server only accepts JSON
func preflightAsset(a *[]byte) error {
	if len(*a) == 0 {
		return errors.New("input must not be empty")
	}

	if utf8.Valid(*a) == false {
		return errors.New("input must be valid UTF-8")
	}

	// attempt to parse JSON first
	var any interface{}
	err := json.Unmarshal(*a, &any)

	// input is valid JSON
	if err == nil {
		return nil
	}

	jsonError := err

	// not JSON
	json, err := yaml.YAMLToJSON(*a)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid JSON: %v; invalid YAML: %v", jsonError, err))
	}

	// successful conversion
	*a = json

	return nil
}
