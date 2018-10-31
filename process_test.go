package main

import (
	"testing"
)

func TestProcessFileInvalidPath(t *testing.T) {
	invalidPath := "/non/existent/file.yaml"
	output := "dot"
	_, err := processFile(invalidPath, &output)

	if err == nil {
		t.Errorf("Must reject invalid path %s", invalidPath)
	}
}

func TestProcessBytes(t *testing.T) {
	//don't allow XML
	xmlBuffer := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="true"?><root/>`)
	output := "dot"
	_, err := processBytes(xmlBuffer, &output)

	if err == nil {
		t.Errorf("Must reject XML input")
	}
}
