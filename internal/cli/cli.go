package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	ExitOK          = 0
	ExitValidation  = 1
	ExitUsage       = 2
	ExitEnvironment = 3
	ExitInternal    = 4
)

type Options struct {
	Repo        string
	RoadmapRoot string
	Workspace   bool
	Output      string
	Strict      bool
	Rootline    string
	Timeout     time.Duration
}

func Execute(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printRootHelp(stdout)
		return ExitOK
	}

	if args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		printRootHelp(stdout)
		return ExitOK
	}

	switch args[0] {
	case "doctor":
		return executeLeafCommand(args[1:], stdout, stderr, "doctor", "Diagnose repo, roadmap configuration and Rootline availability.")
	case "check":
		return executeLeafCommand(args[1:], stdout, stderr, "check", "Validate canonical roadmap structure, metadata and dependency graph.")
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		return ExitUsage
	}
}

func executeLeafCommand(args []string, stdout io.Writer, stderr io.Writer, name string, description string) int {
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help") {
		fs := newFlagSet(name, stdout)
		fs.Usage()
		return ExitOK
	}

	fs := newFlagSet(name, stderr)
	if err := fs.Parse(args); err != nil {
		return ExitUsage
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "%s: unexpected argument %q\n", name, fs.Arg(0))
		return ExitUsage
	}
	fmt.Fprintf(stdout, "%s: not implemented yet\n", name)
	return ExitOK
}

func newFlagSet(name string, output io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(output)
	fs.String("repo", ".", "repository root or workspace member to inspect")
	fs.String("roadmap-root", "", "override configured roadmap root")
	fs.Bool("workspace", false, "treat repo as workspace")
	fs.String("output", "text", "output format: text or json")
	fs.Bool("strict", false, "treat warnings as failures")
	fs.String("rootline", "", "rootline executable path")
	fs.Duration("timeout", 10*time.Second, "timeout for each Rootline call")
	fs.Usage = func() {
		fmt.Fprintf(output, "%s\n\nUsage:\n  roadmapctl %s [flags]\n\nFlags:\n", commandDescription(name), name)
		fs.PrintDefaults()
	}
	return fs
}

func commandDescription(name string) string {
	switch name {
	case "doctor":
		return "Diagnose repo, roadmap configuration and Rootline availability."
	case "check":
		return "Validate canonical roadmap structure, metadata and dependency graph."
	default:
		return ""
	}
}

func printRootHelp(w io.Writer) {
	lines := []string{
		"roadmapctl validates Rootline-governed roadmaps.",
		"",
		"Usage:",
		"  roadmapctl [global flags] <command> [flags]",
		"",
		"Commands:",
		"  doctor    Diagnose repo, roadmap config, Rootline availability and schema prerequisites.",
		"  check     Validate canonical roadmap structure, metadata and dependencies.",
		"",
		"Global flags:",
		"  --repo path             repository root or workspace member to inspect",
		"  --roadmap-root path     override configured roadmap root",
		"  --workspace             treat repo as workspace",
		"  --output text|json      select output format",
		"  --strict                treat warnings as failures",
		"  --rootline path         rootline executable path",
		"  --timeout duration      timeout for each Rootline call",
	}
	fmt.Fprintln(w, strings.Join(lines, "\n"))
}
