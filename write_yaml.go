package main

import (
	"bytes"
	"fmt"

	"github.com/ghodss/yaml"
)

func writeYaml(percentageIsolated int, percentageIsolatedNamespaceToNamespace int, percentageNamespaceCoverage int, buffer *bytes.Buffer) {
	result := Result{percentageIsolated, percentageIsolatedNamespaceToNamespace, percentageNamespaceCoverage}
	yaml, err := yaml.Marshal(&result)
	if err != nil {
		fmt.Printf("Can't encode as YAML: %s", err)
		return
	}

	buffer.Write(yaml)
}
