package main

import (
	"bytes"
	"fmt"
	"strconv"
)

//generate basic markdown table from slice of string slices
func markdownTable(t *[][]string) string {
	maxLengthMap := make(map[int]int)

	//1st pass: populate maxLengthMap
	for rowIndex, row := range *t {
		for cellIndex, cell := range row {
			length := len(cell)
			if rowIndex == 0 {
				length += 4
			}
			if length > maxLengthMap[cellIndex] {
				maxLengthMap[cellIndex] = length
			}
		}
	}

	var b bytes.Buffer
	for rowIndex, row := range *t {
		for cellIndex, cell := range row {
			if rowIndex == 0 {
				cell = "**" + cell + "**"
			}
			if cellIndex == 0 {
				b.WriteByte('|')
			}
			b.WriteString(fmt.Sprintf("%-"+strconv.Itoa(maxLengthMap[cellIndex])+"s", cell))
			b.WriteByte('|')
		}
		b.WriteByte('\n')
		if rowIndex == 0 {
			for cellIndex, _ := range row {
				if cellIndex == 0 {
					b.WriteByte('|')
				}
				b.WriteByte(':')
				for i := 0; i < (maxLengthMap[cellIndex] - 1); i++ {
					b.WriteByte('-')
				}
				b.WriteByte('|')
			}
			b.WriteByte('\n')
		}
	}
	b.WriteByte('\n')
	return b.String()
}
