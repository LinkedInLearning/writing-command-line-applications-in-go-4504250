package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/353solutions/infra/logs"
)

var parseConfig struct {
	count   int
	verbose bool
}

var parseUsage = `Usage: %s parse [options] [FILE]
Validate server logs in FILE (standard input if not given).

Options:
`

func runParse(args []string) {
	if err := parseEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fs := flag.NewFlagSet("parse", flag.ExitOnError)
	fs.Usage = func() {
		name := path.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, parseUsage, name)
		fs.PrintDefaults()
	}
	fs.Var(&Count{&parseConfig.count}, "count", "number of records to parse")
	fs.BoolVar(&parseConfig.verbose, "verbose", parseConfig.verbose, "emit more information (also LOGS_VERBOSE)")
	fs.Parse(args[1:])

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
		if parseConfig.count > 0 && parseConfig.count == n {
			break
		}
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s:%d: scanning - %s\n", fileName, s.Line(), err)
		os.Exit(1)
	}

	if parseConfig.verbose {
		fmt.Printf("%s: successfully processed %d records\n", fileName, n)
	}
}

func parseEnv() error {
	const verboseKey = "LOGS_VERBOSE"
	v := os.Getenv(verboseKey)
	switch v {
	// See https://pkg.go.dev/flag#hdr-Command_line_flag_syntax
	case "1", "t", "T", "true", "TRUE", "True":
		parseConfig.verbose = true
	case "0", "f", "F", "false", "FALSE", "False":
		parseConfig.verbose = false
	case "":
		// NOP
	default:
		return fmt.Errorf("bad value for %s - %q\n", verboseKey, v)
	}

	return nil
}
