package main

import (
	"bytes"
	"fmt"
)

func writeMarkdown(percentageIsolatedInt int, percentageNamespaceCoverageInt int, b *bytes.Buffer) {
	heading("Network policy viewer", 1, b)

	rows := make([][]string, 3) //set length plus header row
	rows[0] = append(rows[0], "Field")
	rows[0] = append(rows[0], "Value")
	rows[1] = append(rows[1], "percentageIsolated")
	rows[1] = append(rows[1], fmt.Sprintf("%d", percentageIsolatedInt))
	rows[2] = append(rows[2], "percentageNamespaceCoverage")
	rows[2] = append(rows[2], fmt.Sprintf("%d", percentageNamespaceCoverageInt))

	b.WriteString(markdownTable(&rows))
}

func heading(s string, level int, b *bytes.Buffer) {
	if level < 3 {
		b.WriteString(s + "\n")
		for i := 0; i < len(s); i++ {
			underline := "="
			if level > 1 {
				underline = "-"
			}
			b.WriteString(underline)
		}
	} else {
		for i := 0; i < level; i++ {
			b.WriteString("#")
		}
		b.WriteString(" " + s)
	}
	b.WriteString("\n\n")
}
