package main

import (
	"strings"
	"testing"
)

func TestPage(t *testing.T) {
	unescapedHTML := "<p>Alice, Bob &amp; Eve</p>"
	result := page("title", unescapedHTML)
	if !strings.Contains(result, unescapedHTML) {
		t.Errorf("Unescaped HTML has to be passed unchanged")
	}
}
