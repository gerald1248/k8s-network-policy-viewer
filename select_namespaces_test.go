package main

import (
	"testing"
)

func TestMatchingNamespace(t *testing.T) {
	selector := map[string]string{}
	selector["app"] = "alice"
	namespacePodMap := map[string][]string{}
	nestedArray := []string{}
	nestedArray = append(nestedArray, "alice")
	namespacePodMap["alice"] = nestedArray
	namespaceLabelMap := map[string]map[string]string{}
	nestedSelector := map[string]string{}
	nestedSelector["app"] = "alice"
	namespaceLabelMap["alice"] = nestedSelector

	namespaces := selectNamespaces(&selector, &namespacePodMap, &namespaceLabelMap)

	if len(namespaces) != 1 || namespaces[0] != "alice" {
		t.Errorf("Must find matching namespace - got array=%v", namespaces)
	}
}
