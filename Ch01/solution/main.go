package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

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

var usage = `Usage: %s <options> [TEXT]
Prints a banner. For example:
$ banner -w 6 Go
  Go
======
Options:
`

func main() {
	width := 80
	var out io.Writer = os.Stdout
	outFile := ""

	flag.Var(&Width{&width}, "width", "banner width")
	flag.StringVar(&outFile, "out", "", "output file (default stdout)")
	flag.Usage = func() {
		name := path.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, usage, name)
		flag.PrintDefaults()
	}
	flag.Parse()

	text := ""

	switch flag.NArg() {
	case 0:
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: can't read - %s\n", err)
			os.Exit(1)
		}
		text = string(data)
	case 1:
		text = flag.Arg(0)
	default:
		fmt.Fprintln(os.Stderr, "error: wrong number of arguments")
		os.Exit(1)
	}

	if outFile != "" {
		file, err := os.Create(outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		out = file
	}

	Banner(out, text, width)
}

type Width struct {
	w *int
}

func (w *Width) String() string {
	if w.w == nil {
		return ""
	}

	return fmt.Sprintf("%d", *w.w)
}

func (w *Width) Set(val string) error {
	i, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	if i <= 0 || i > 250 {
		return fmt.Errorf("width %d out of range [1:250]", i)
	}

	*w.w = i
	return nil
}
