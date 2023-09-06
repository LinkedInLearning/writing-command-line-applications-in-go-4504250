package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/353solutions/infra/logs"
)

var countUsage = `Usage: %s count [options] [FILE]
Count server logs in FILE (standard input if not given).

Options:
`

func runCount(args []string) {
	if err := parseEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fs := flag.NewFlagSet("count", flag.ExitOnError)
	fs.Usage = func() {
		name := path.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, countUsage, name)
		fs.PrintDefaults()
	}
	fs.Parse(args[1:])

	// TODO: Unite with validate.go
	var r io.Reader
	var fileName string
	switch fs.NArg() {
	case 0:
		r, fileName = os.Stdin, "<stdin>"
	case 1:
		fileName = fs.Arg(0)
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()

		r = file
	default:
		fmt.Fprintln(os.Stderr, "error: wrong number of arguments")
		os.Exit(1)
	}

	s := logs.NewScanner(r)
	n := 0
	for s.Next() {
		n++
		if validateCfg.count > 0 && validateCfg.count == n {
			break
		}
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s:%d: scanning - %s\n", fileName, s.Line(), err)
		os.Exit(1)
	}

	fmt.Println(n)
}
