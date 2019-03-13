package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// make individual config items part of a list, so only one type is required for the call to json.Unmarshal
func makeList(a *[]byte) error {
	var obj MinimalObject
	err := json.Unmarshal(*a, &obj)

	if err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	switch obj.Kind {
	case "List", "Template":
		return nil
	case "Pod", "NetworkPolicy", "Ingress":
		slices := [][]byte{[]byte(`{"kind":"List","items":[`), *a, []byte(`]}`)}
		b := bytes.Join(slices, []byte{})
		*a = b
		return nil
	}
	return fmt.Errorf("can't parse JSON: no API objects found")
}
