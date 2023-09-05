package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/353solutions/infra/logs"
)

var config struct {
	count   int
	verbose bool
}

var usage = `Usage: %s [options] [FILE]
Validate server logs in FILE (standard input if not given).

Options:
`

func main() {
	if err := parseEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	flag.Usage = func() {
		name := path.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, usage, name)
		flag.PrintDefaults()
	}
	flag.Var(&Count{&config.count}, "count", "number of records to parse")
	flag.BoolVar(&config.verbose, "verbose", config.verbose, "emit more information (also LOGS_VERBOSE)")
	flag.Parse()

	var r io.Reader
	var fileName string
	switch flag.NArg() {
	case 0:
		r, fileName = os.Stdin, "<stdin>"
	case 1:
		fileName = flag.Arg(0)
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
		if config.count > 0 && config.count == n {
			break
		}
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s:%d: scanning - %s\n", fileName, s.Line(), err)
		os.Exit(1)
	}

	if config.verbose {
		fmt.Printf("%s: successfully validated %d records\n", fileName, n)
	}
}

func parseEnv() error {
	const verboseKey = "LOGS_VERBOSE"
	v := os.Getenv(verboseKey)
	switch v {
	// See https://pkg.go.dev/flag#hdr-Command_line_flag_syntax
	case "1", "t", "T", "true", "TRUE", "True":
		config.verbose = true
	case "0", "f", "F", "false", "FALSE", "False":
		config.verbose = false
	case "":
		// NOP
	default:
		return fmt.Errorf("bad value for %s - %q\n", verboseKey, v)
	}

	return nil
}
