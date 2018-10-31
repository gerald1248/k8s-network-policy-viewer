package main

import "bytes"
import "fmt"

func writeDot(namespacePodMap *map[string][]string, buffer *bytes.Buffer) {
	buffer.WriteString("digraph podNetwork {\n")
	counter := 0
	var allPods []string
	for k, v := range *namespacePodMap {
		counter += 1
		buffer.WriteString("  subgraph cluster_")
		fmt.Fprintf(buffer, "%d", counter)
		buffer.WriteString(" {\n")
		for _, s := range v {
			buffer.WriteString("    \"")
			buffer.WriteString(s)
			buffer.WriteString("\";\n")
			allPods = append(allPods, s)
		}
		buffer.WriteString("    label = \"")
		buffer.WriteString(k)
		buffer.WriteString("\"\n")
		buffer.WriteString("  }\n")
	}
	// connect all for now
	for _, outer := range allPods {
		for _, inner := range allPods {
			fmt.Fprintf(buffer, "  \"%s\" -> \"%s\";\n", outer, inner)
		}
	}
	buffer.WriteString("}\n")
}
