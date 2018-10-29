package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// make individual config items part of a list, so only one type is required for the call to json.Unmarshal
func makeList(a *[]byte) error {
	var obj MinimalObject
	err := json.Unmarshal(*a, &obj)

	if err != nil {
		return errors.New(fmt.Sprintf("invalid JSON: %v", err))
	}

	switch obj.Kind {
	case "List", "Template":
		return nil
	case "BuildConfig", "DeploymentConfig", "Pod", "Route":
		slices := [][]byte{[]byte(`{"kind":"List","items":[`), *a, []byte(`]}`)}
		b := bytes.Join(slices, []byte{})
		*a = b
		return nil
	}
	return errors.New(fmt.Sprintf("can't parse JSON: no configuration object found"))
}
