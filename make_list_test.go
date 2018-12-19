package main

import "testing"

func TestMakeList(t *testing.T) {
	jsonBuffer := []byte(`{"kind":"Pod","metadata":{},"spec":{}}`)
	err := makeList(&jsonBuffer)

	if err != nil {
		t.Errorf("Can't process byte array: %s\n", err)
	}

	expected := `{"kind":"List","items":[{"kind":"Pod","metadata":{},"spec":{}}]}`
	actual := string(jsonBuffer)
	if actual != expected {
		t.Errorf("Unexpected result\nExpected:\n%s\nGot:\n%s\n", expected, actual)
	}
}
