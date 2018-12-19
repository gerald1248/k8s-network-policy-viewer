package main

import (
	"strings"
	"testing"
)

func TestPage(t *testing.T) {
	unescapedHtml := "<p>Alice, Bob &amp; Eve</p>"
	result := page(unescapedHtml)
	if !strings.Contains(result, unescapedHtml) {
		t.Errorf("Unescaped HTML has to be passed unchanged")
	}
}
