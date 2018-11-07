package main

import "bytes"
import "encoding/json"
import "fmt"

func writeJson(namespacePodMap *map[string][]string, buffer *bytes.Buffer) {
	json, err := json.Marshal(namespacePodMap)
	if err != nil {
		fmt.Printf("Can't encode as JSON: %s", err)
		return
	}

	buffer.Write(json)
}
