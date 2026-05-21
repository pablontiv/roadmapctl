package main

import (
	"io"
	"os"

	"github.com/pablontiv/roadmapctl/internal/cli"
)

// mainFn is a variable that points to the actual os.Exit,
// allowing it to be mocked in tests.
var mainFn = os.Exit

func main() {
	mainFn(run(os.Args[1:], os.Stdout, os.Stderr, version))
}

func run(args []string, stdout io.Writer, stderr io.Writer, version string) int {
	return cli.Execute(args, stdout, stderr, version)
}
