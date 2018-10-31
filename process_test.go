package main

import (
	"testing"
)

func TestProcessFileInvalidPath(t *testing.T) {
	invalidPath := "/non/existent/file.yaml"
	_, err := processFile(invalidPath)

	if err == nil {
		t.Errorf("Must reject invalid path %s", invalidPath)
	}
}

func TestProcessBytes(t *testing.T) {
	//don't allow XML
	xmlBuffer := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="true"?><root/>`)
	_, err := processBytes(xmlBuffer)

	if err == nil {
		t.Errorf("Must reject XML input")
	}
}
