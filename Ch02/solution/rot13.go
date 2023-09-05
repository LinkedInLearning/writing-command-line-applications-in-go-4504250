package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
)

var usage = `Usage: %s <options> [TEXT]
Encode text in rot13 encoding.
Options:
`

func main() {
	jsonFormat := false
	flag.BoolVar(&jsonFormat, "json", false, "emit output in JSON format")
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

	encoded := rot13(text)
	if !jsonFormat {
		fmt.Println(encoded)
		return
	}

	data := map[string]any{
		"text":    text,
		"encoded": encoded,
	}
	json.NewEncoder(os.Stdout).Encode(data)
}

func rot13(text string) string {
	out := make([]byte, len(text))
	for i := 0; i < len(text); i++ {
		c := text[i]
		switch {
		case (c >= 'a' && c <= 'm') || (c >= 'A' && c <= 'M'):
			c += 13
		case (c >= 'n' && c <= 'z') || (c >= 'N' && c <= 'Z'):
			c -= 13
		}
		out[i] = c
	}

	return string(out)
}
