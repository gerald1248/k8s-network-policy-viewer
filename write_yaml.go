package main

import (
	"bytes"
	"fmt"

	"github.com/ghodss/yaml"
)

func writeYaml(percentageIsolatedInt int, percentageNamespaceCoverageInt int, buffer *bytes.Buffer) {
	result := Result{percentageIsolatedInt, percentageNamespaceCoverageInt}
	yaml, err := yaml.Marshal(&result)
	if err != nil {
		fmt.Printf("Can't encode as YAML: %s", err)
		return
	}

	buffer.Write(yaml)
}
