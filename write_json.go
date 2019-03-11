package main

import "bytes"
import "encoding/json"
import "fmt"

func writeJson(percentageIsolated int, percentageIsolatedNamespaceToNamespace int, percentageNamespaceCoverage int, buffer *bytes.Buffer) {
	result := Result{percentageIsolated, percentageIsolatedNamespaceToNamespace, percentageNamespaceCoverage}
	json, err := json.Marshal(&result)
	if err != nil {
		fmt.Printf("Can't encode as JSON: %s", err)
		return
	}

	buffer.Write(json)
}
