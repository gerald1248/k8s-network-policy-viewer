package main

import (
	"bytes"
	"testing"
)

func TestHeading(t *testing.T) {
	expectedH1 := `Foobar
======

`
	expectedH3 := `### Foobar

`

	var bufferH1 bytes.Buffer
	heading("Foobar", 1, &bufferH1)
	if bufferH1.String() != expectedH1 {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expectedH1, bufferH1.String())
	}

	var bufferH3 bytes.Buffer
	heading("Foobar", 3, &bufferH3)
	if bufferH3.String() != expectedH3 {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expectedH3, bufferH3.String())
	}
}
