package main

import (
	"bytes"
	"testing"
)

func TestWriteYaml(t *testing.T) {
	expected := `percentageIsolated: 20
percentageNamespaceCoverage: 30
`

	var buffer bytes.Buffer
	writeYaml(20, 30, &buffer)
	s := buffer.String()
	if s != expected {
		t.Errorf("YAML output faulty\nExpected:\n%s\nGot:\n%s\n", expected, s)
	}
}
