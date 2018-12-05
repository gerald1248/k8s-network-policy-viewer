package main

import (
	"bytes"
	"fmt"
	"strings"
)

func writeDot(namespacePodMap *map[string][]string, edgeMap *map[string][]string, buffer *bytes.Buffer) {
	buffer.WriteString("digraph podNetwork {\n")
	counter := 0
	var allPods []string
	for k, v := range *namespacePodMap {
		counter += 1
		buffer.WriteString("  subgraph cluster_")
		fmt.Fprintf(buffer, "%d", counter)
		buffer.WriteString(" {\n")
		for _, s := range v {
			sUnqualified := stripNamespace(s)
			buffer.WriteString("    \"")
			buffer.WriteString(sUnqualified)
			buffer.WriteString("\";\n")
			allPods = append(allPods, sUnqualified)
		}
		buffer.WriteString("    label = \"")
		buffer.WriteString(k)
		buffer.WriteString("\"\n")
		buffer.WriteString("  }\n")
	}
	for k, v := range *edgeMap {
		for _, s := range v {
			kUnqualified := stripNamespace(k)
			sUnqualified := stripNamespace(s)
			fmt.Fprintf(buffer, "  \"%s\" -> \"%s\";\n", kUnqualified, sUnqualified)
		}
	}
	buffer.WriteString("}\n")
}

func stripNamespace(s string) string {
	index := strings.IndexRune(s, ':')
	if index == -1 {
		return s
	}
	return s[index+1 : len(s)]
}
