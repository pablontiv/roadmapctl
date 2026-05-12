package main

import (
	"io"
	"os"

	"github.com/pablontiv/roadmapctl/internal/cli"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, version))
}

func run(args []string, stdout io.Writer, stderr io.Writer, version string) int {
	return cli.Execute(args, stdout, stderr, version)
}
