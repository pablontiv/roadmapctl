package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

const (
	ExitOK          = diagnostics.ExitOK
	ExitValidation  = diagnostics.ExitValidation
	ExitUsage       = diagnostics.ExitUsage
	ExitEnvironment = diagnostics.ExitEnvironment
	ExitInternal    = diagnostics.ExitInternal
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
	options := Options{Repo: ".", Output: "text", Timeout: 10 * time.Second}
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help") {
		fs := newFlagSet(name, stdout, &options)
		fs.Usage()
		return ExitOK
	}

	fs := newFlagSet(name, stderr, &options)
	if err := fs.Parse(args); err != nil {
		return ExitUsage
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "%s: unexpected argument %q\n", name, fs.Arg(0))
		return ExitUsage
	}
	if options.Output != "text" && options.Output != "json" {
		fmt.Fprintf(stderr, "%s: unsupported output format %q\n", name, options.Output)
		return ExitUsage
	}

	report := diagnostics.NewReport("roadmapctl/"+name, options.Repo, options.RoadmapRoot, nil)
	if options.Output == "json" {
		if err := diagnostics.RenderJSON(stdout, report); err != nil {
			fmt.Fprintf(stderr, "%s: render JSON report: %v\n", name, err)
			return ExitInternal
		}
		return diagnostics.ExitCode(report, options.Strict)
	}
	if err := diagnostics.RenderText(stdout, report); err != nil {
		fmt.Fprintf(stderr, "%s: render text report: %v\n", name, err)
		return ExitInternal
	}
	return diagnostics.ExitCode(report, options.Strict)
}

func newFlagSet(name string, output io.Writer, options *Options) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(output)
	fs.StringVar(&options.Repo, "repo", options.Repo, "repository root or workspace member to inspect")
	fs.StringVar(&options.RoadmapRoot, "roadmap-root", options.RoadmapRoot, "override configured roadmap root")
	fs.BoolVar(&options.Workspace, "workspace", options.Workspace, "treat repo as workspace")
	fs.StringVar(&options.Output, "output", options.Output, "output format: text or json")
	fs.BoolVar(&options.Strict, "strict", options.Strict, "treat warnings as failures")
	fs.StringVar(&options.Rootline, "rootline", options.Rootline, "rootline executable path")
	fs.DurationVar(&options.Timeout, "timeout", options.Timeout, "timeout for each Rootline call")
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
