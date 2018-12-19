package main

import "bytes"
import "encoding/json"
import "fmt"

func writeJson(percentageIsolatedInt int, percentageNamespaceCoverageInt int, buffer *bytes.Buffer) {
	result := Result{percentageIsolatedInt, percentageNamespaceCoverageInt}
	json, err := json.Marshal(&result)
	if err != nil {
		fmt.Printf("Can't encode as JSON: %s", err)
		return
	}

	buffer.Write(json)
}
