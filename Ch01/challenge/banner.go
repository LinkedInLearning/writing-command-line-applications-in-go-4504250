package main

import (
	"fmt"
	"io"
	"unicode"
)

// Banners print text centered to width and then "===" under it.
func Banner(w io.Writer, text string, width int) {
	n := numPrintable(text)
	padding := (width - n) / 2
	fmt.Fprintf(w, "%*s\n", padding+n, text)
	for i := 0; i < width; i++ {
		w.Write([]byte{'='})
	}
	w.Write([]byte{'\n'})
}

// numPrintable returns the number of printable runes in text
func numPrintable(text string) int {
	count := 0
	for _, r := range text {
		if unicode.IsPrint(r) {
			count++
		}
	}
	return count
}
