package main

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

var (
	commands = make(map[string]func([]string))
)

var usage = `Usage: %s <command> [<args>]
Where command is one of: %s
Run "%s <command> -h" for sub command help.
`

func main() {
	commands["validate"] = runValidate
	commands["count"] = runCount

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: missing command\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "-h", "--help":
		name := path.Base(os.Args[0])
		cmds := strings.Join(allCommands(), "|")
		fmt.Fprintf(os.Stderr, usage, name, cmds, name)
		os.Exit(0)
	}

	name := os.Args[1]
	fn, ok := commands[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "error: unknown command - %q\n", name)
		os.Exit(1)
	}
	fn(os.Args[1:])
}

func allCommands() []string {
	cmds := make([]string, 0, len(commands))
	for name := range commands {
		cmds = append(cmds, name)
	}
	sort.Strings(cmds)
	return cmds
}
