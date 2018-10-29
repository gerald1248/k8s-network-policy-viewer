package main

import "bytes"

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

func table(s ContainerSet, b *bytes.Buffer) {
	//alloc string table
	count := len(s)
	rows := make([][]string, count+1) //set length plus header row

	//header row
	rows[0] = []string{"Namespace", "Name", "Container"}

	//special case: one empty spec
	if count == 1 && s[0].Namespace == "" && s[0].Name == "" && s[0].Container == "" {
		//exit without writing to buffer
		return
	}

	//content rows
	for i := 0; i < count; i++ {
		spec := s[i]
		rows[i+1] = []string{spec.Namespace, spec.Name, spec.Container}
	}

	b.WriteString(markdownTable(&rows))
}
