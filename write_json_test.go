package main

import (
	"bytes"
	"testing"
)

func TestWriteJson(t *testing.T) {
	expected := `{"percentageIsolated":20,"percentageNamespaceCoverage":30}`

	var buffer bytes.Buffer
	writeJson(20, 30, &buffer)
	s := buffer.String()
	if s != expected {
		t.Errorf("JSON output faulty\nExpected:\n%s\nGot:\n%s\n", expected, s)
	}
}
