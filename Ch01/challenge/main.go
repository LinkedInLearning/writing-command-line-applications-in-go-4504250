package main

import "os"

/*
Instructions:

- Pass text to print as an argument
  - If not argument - read from stdin

- Use -width to specify width
  - Width should be bigger than 0 and less than 250
  - Default to 80

- Use -out to specify output file
  - Default to stdout
*/
func main() {
	text := "Go"
	width := 6
	Banner(os.Stdout, text, width)
}
