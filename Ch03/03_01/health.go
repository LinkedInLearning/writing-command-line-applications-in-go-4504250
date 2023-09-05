package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	var timeout time.Duration
	var verbose bool

	flag.DurationVar(&timeout, "timeout", 3*time.Second, "request timeout")
	flag.BoolVar(&verbose, "verbose", false, "emit more information")
	flag.Usage = func() {
		name := path.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [options] URL\n", name)
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "error: wrong number of arguments\n")
		os.Exit(1)
	}

	url := flag.Arg(0)
	if verbose {
		fmt.Printf("checking %q\n", url)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: can't create request: %s\n", err)
		os.Exit(1)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %q: can't call - %s\n", url, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error: %q: bad status - %s\n", url, resp.Status)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("%q: %s\n", url, resp.Status)
	}
}
