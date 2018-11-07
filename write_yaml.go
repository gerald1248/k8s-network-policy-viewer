package main

import "bytes"
import "fmt"
import "github.com/ghodss/yaml"

func writeYaml(namespacePodMap *map[string][]string, buffer *bytes.Buffer) {
	yaml, err := yaml.Marshal(namespacePodMap)
	if err != nil {
		fmt.Printf("Can't encode as YAML: %s", err)
		return
	}

	buffer.Write(yaml)
}
