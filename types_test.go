package main

import "testing"

func TestMetadata(t *testing.T) {
	labels := make(map[string]string)
	annotations := make(map[string]string)
	nameValue := "name"
	namespaceValue := "namespace"
	metadata := Metadata{labels, annotations, nameValue, namespaceValue}
	if len(metadata.Labels) > 0 || metadata.Name != nameValue || metadata.Namespace != namespaceValue {
		t.Errorf("faulty struct metadata: %v", metadata)
	}
}
