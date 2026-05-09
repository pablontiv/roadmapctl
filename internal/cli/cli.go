package cli

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/spf13/cobra"
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
	exitCode := ExitOK
	cmd := newRootCommand(stdout, stderr, &exitCode)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitUsage
	}
	return exitCode
}

func newRootCommand(stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	options := Options{Repo: ".", Output: "text", Timeout: 10 * time.Second}
	cmd := &cobra.Command{
		Use:           "roadmapctl",
		Short:         "roadmapctl validates Rootline-governed roadmaps.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.CompletionOptions.DisableDefaultCmd = true

	flags := cmd.PersistentFlags()
	flags.StringVar(&options.Repo, "repo", options.Repo, "repository root or workspace member to inspect")
	flags.StringVar(&options.RoadmapRoot, "roadmap-root", options.RoadmapRoot, "override configured roadmap root")
	flags.BoolVar(&options.Workspace, "workspace", options.Workspace, "treat repo as workspace")
	flags.StringVar(&options.Output, "output", options.Output, "output format: text or json")
	flags.BoolVar(&options.Strict, "strict", options.Strict, "treat warnings as failures")
	flags.StringVar(&options.Rootline, "rootline", options.Rootline, "rootline executable path")
	flags.DurationVar(&options.Timeout, "timeout", options.Timeout, "timeout for each Rootline call")

	cmd.AddCommand(newLeafCommand("doctor", "Diagnose repo, roadmap configuration and Rootline availability.", &options, stdout, stderr, exitCode))
	cmd.AddCommand(newLeafCommand("check", "Validate canonical roadmap structure, metadata and dependency graph.", &options, stdout, stderr, exitCode))
	cmd.AddCommand(newLeafCommand("context", "Show effective roadmapctl context for skill bootstrap.", &options, stdout, stderr, exitCode))
	cmd.AddCommand(newLeafCommand("pending", "List active roadmap tasks that are not done.", &options, stdout, stderr, exitCode))
	cmd.AddCommand(newBootstrapCommand(&options, stdout, stderr, exitCode))
	return cmd
}

func newLeafCommand(name string, description string, options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	return &cobra.Command{
		Use:           name,
		Short:         description,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			*exitCode = executeLeafCommand(context.Background(), name, *options, stdout, stderr)
			return nil
		},
	}
}

func executeLeafCommand(ctx context.Context, name string, options Options, stdout io.Writer, stderr io.Writer) int {
	if options.Output != "text" && options.Output != "json" {
		fmt.Fprintf(stderr, "%s: unsupported output format %q\n", name, options.Output)
		return ExitUsage
	}

	if name == "pending" {
		report := runPending(ctx, options)
		if options.Output == "json" {
			if err := renderPendingJSON(stdout, report); err != nil {
				fmt.Fprintf(stderr, "%s: render JSON report: %v\n", name, err)
				return ExitInternal
			}
			return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), options.Strict)
		}
		if err := renderPendingText(stdout, report); err != nil {
			fmt.Fprintf(stderr, "%s: render text report: %v\n", name, err)
			return ExitInternal
		}
		return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), options.Strict)
	}

	if name == "context" {
		report := runContext(ctx, options)
		if options.Output == "json" {
			if err := renderContextJSON(stdout, report); err != nil {
				fmt.Fprintf(stderr, "%s: render JSON report: %v\n", name, err)
				return ExitInternal
			}
			return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), options.Strict)
		}
		if err := renderContextText(stdout, report); err != nil {
			fmt.Fprintf(stderr, "%s: render text report: %v\n", name, err)
			return ExitInternal
		}
		return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), options.Strict)
	}

	report := diagnostics.NewReport("roadmapctl/"+name, options.Repo, options.RoadmapRoot, nil)
	if name == "doctor" {
		report = runDoctor(ctx, options)
	}
	if name == "check" {
		report = runCheck(ctx, options)
	}
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
