package main

import (
	"testing"
)

func TestPreflightAsset(t *testing.T) {
	//byte slices
	invalidUtf8 := []byte{0xff, 0xfe, 0xfd}
	validJSON := []byte("{ \"foo\": [\"bar\", \"barfoo\"] }")
	invalidJSON := []byte("foo { \"foo\": [\"bar\", \"barfoo\"] } foo")
	validYAML := []byte("\"foo\": \"bar\"")
	invalidYAML := []byte("\"foo\":foo:bar: \"bar\"")
	multilineYAML := []byte(`"foo":
- "bar"
- "foobar"
- "boofar"
- "roobar"
`)
	multilineYAMLConverted := []byte("{\"foo\":[\"bar\",\"foobar\",\"boofar\",\"roobar\"]}")

	//expect error
	err := preflightAsset(&invalidUtf8)
	if err == nil {
		t.Error("Must reject invalid UTF8")
	}

	err = preflightAsset(&invalidJSON)
	if err == nil {
		t.Error("Must reject invalid JSON")
	}

	err = preflightAsset(&invalidYAML)
	if err == nil {
		t.Error("Must reject invalid YAML")
	}

	//expect success
	err = preflightAsset(&validYAML)
	if err != nil {
		t.Errorf("Must accept valid YAML: %v", err)
	}

	err = preflightAsset(&validJSON)
	if err != nil {
		t.Errorf("Must accept valid JSON: %v", err)
	}

	//in-place conversion must match predefined result
	err = preflightAsset(&multilineYAML)
	if err != nil {
		t.Errorf("Must accept valid multiline YAML: %v", err)
	}
	if string(multilineYAML) != string(multilineYAMLConverted) {
		t.Errorf("Expected %s to match %s", multilineYAML, multilineYAMLConverted)
	}
}
