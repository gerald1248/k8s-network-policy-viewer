package main

import (
	"testing"
)

func TestMarkdownTable(t *testing.T) {
	rows := make([][]string, 2)
	rows[0] = []string{"A", "B"}
	rows[1] = []string{"C", "D"}

	//add padding
	expected := `|**A**|**B**|
|:----|:----|
|C    |D    |

`

	md := markdownTable(&rows)

	if md != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expected, md)
	}
}
