package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/roadmap"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
	"github.com/spf13/cobra"
)

type transitionReport struct {
	Version              int                          `json:"version"`
	Kind                 string                       `json:"kind"`
	Summary              diagnostics.Summary          `json:"summary"`
	Root                 string                       `json:"root"`
	RoadmapRoot          string                       `json:"roadmap_root"`
	Action               string                       `json:"action"`
	Path                 string                       `json:"path"`
	Allowed              bool                         `json:"allowed"`
	CurrentStatus        string                       `json:"current_status"`
	TargetStatus         string                       `json:"target_status"`
	Role                 string                       `json:"role"`
	Reasons              []string                     `json:"reasons"`
	BlockingDependencies []roadmap.BlockingDependency `json:"blocking_dependencies"`
	Changes              []roadmap.TransitionChange   `json:"changes"`
	Diagnostics          []diagnostics.Diagnostic     `json:"diagnostics"`
}

func newTransitionCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	cmd := &cobra.Command{Use: "transition", Short: "Evaluate policy-checked roadmap status transitions.", SilenceUsage: true, SilenceErrors: true}
	for _, action := range []string{"can-start", "can-complete", "start", "complete"} {
		action := action
		cmd.AddCommand(newTransitionActionCommand(options, stdout, stderr, exitCode, action, ""))
	}
	setStatus := newTransitionActionCommand(options, stdout, stderr, exitCode, "set-status", "")
	setStatus.Flags().String("status", "", "target status value or role name")
	cmd.AddCommand(setStatus)
	return cmd
}

func newTransitionActionCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int, action string, status string) *cobra.Command {
	var dryRun bool
	var apply bool
	command := &cobra.Command{Use: action + " <task-path>", Short: "Evaluate " + action + " for a roadmap task.", Args: cobra.ExactArgs(1), SilenceUsage: true, SilenceErrors: true, RunE: func(cmd *cobra.Command, args []string) error {
		if action == "set-status" {
			flagStatus, _ := cmd.Flags().GetString("status")
			status = flagStatus
		}
		if !apply && (action == "start" || action == "complete") {
			report := newTransitionReport(absoluteClean(options.Repo), "", action, normalizeTransitionPath("", args[0]), roadmap.TransitionResult{Diagnostics: []diagnostics.Diagnostic{{ID: diagnostics.DiagnosticTransitionApplyFailed, Severity: diagnostics.SeverityError, Message: "transition start and complete require --apply flag", Path: normalizeTransitionPath("", args[0]), ExitCode: diagnostics.ExitUsage}}})
			if options.Output == "json" {
				_ = json.NewEncoder(stdout).Encode(report)
			} else {
				fmt.Fprintf(stdout, "%s\nstatus: %s\naction: %s\npath: %s\nallowed: %t\n", report.Kind, report.Summary.Status, report.Action, report.Path, report.Allowed)
			}
			*exitCode = diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), options.Strict)
			return nil
		}
		report := runTransition(context.Background(), *options, action, args[0], status, apply)
		if options.Output == "json" {
			if err := json.NewEncoder(stdout).Encode(report); err != nil {
				fmt.Fprintf(stderr, "transition: render JSON report: %v\n", err)
				*exitCode = ExitInternal
				return nil
			}
		} else {
			fmt.Fprintf(stdout, "%s\nstatus: %s\naction: %s\npath: %s\nallowed: %t\n", report.Kind, report.Summary.Status, report.Action, report.Path, report.Allowed)
		}
		*exitCode = diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), options.Strict)
		return nil
	}}
	command.Flags().BoolVar(&dryRun, "dry-run", true, "plan transition without applying changes")
	command.Flags().BoolVar(&apply, "apply", false, "apply transition after planning and run postcheck")
	return command
}

func runTransition(ctx context.Context, options Options, action string, taskPath string, explicitStatus string, apply bool) transitionReport {
	repoRoot := absoluteClean(options.Repo)
	cfg, err := config.Load(options.Repo, config.Options{RoadmapRoot: options.RoadmapRoot})
	if err != nil {
		found := []diagnostics.Diagnostic{configDiagnostic(repoRoot, err)}
		return newTransitionReport(repoRoot, "", action, normalizeTransitionPath("", taskPath), roadmap.TransitionResult{Diagnostics: found})
	}
	path := normalizeTransitionPath(cfg.RoadmapRoot, taskPath)
	if strings.HasPrefix(path, filepath.ToSlash(cfg.RoadmapRootRel)+"/") {
		path = strings.TrimPrefix(path, filepath.ToSlash(cfg.RoadmapRootRel)+"/")
	}
	model, found := readModelForConfig(ctx, cfg, options)
	roles := roadmap.TransitionRoles{DoneStatuses: cfg.DoneStatuses, ActiveStatuses: cfg.ActiveStatuses, InProgressStatus: cfg.StatusValues.InProgress, CompletedStatus: cfg.StatusValues.Completed}
	var result roadmap.TransitionResult
	switch action {
	case "can-start", "start":
		result = roadmap.CanStart(model, roles, path)
	case "can-complete", "complete":
		result = roadmap.CanComplete(model, roles, path)
	case "set-status":
		targetStatus := transitionStatusValue(cfg, explicitStatus)
		if targetStatus == "" {
			result = roadmap.TransitionResult{Diagnostics: []diagnostics.Diagnostic{{ID: diagnostics.DiagnosticTransitionStatusUnknown, Severity: diagnostics.SeverityError, Message: "target status is required", Path: path, ExitCode: diagnostics.ExitUsage}}}
		} else {
			result = roadmap.SetStatus(model, roles, path, targetStatus)
		}
	default:
		result = roadmap.TransitionResult{Diagnostics: []diagnostics.Diagnostic{{ID: "RMC_TRANSITION_ACTION_UNKNOWN", Severity: diagnostics.SeverityError, Message: "unsupported transition action", Path: path, ExitCode: diagnostics.ExitUsage}}}
	}
	result.Diagnostics = append(found, result.Diagnostics...)
	result = validateTransitionTargetStatus(ctx, cfg, options, result, path)
	if apply && result.Allowed && len(result.Changes) > 0 {
		result = applyTransitionChanges(ctx, cfg, options, result)
	}
	return newTransitionReport(cfg.RepoRoot, cfg.RoadmapRoot, action, path, result)
}

func applyTransitionChanges(ctx context.Context, cfg *config.Config, options Options, result roadmap.TransitionResult) roadmap.TransitionResult {
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: cfg.RepoRoot, Timeout: options.Timeout})
	for i := range result.Changes {
		change := &result.Changes[i]
		_, err := client.Set(ctx, filepath.Join(cfg.RoadmapRoot, filepath.FromSlash(change.Path)), change.Field+"="+change.After)
		if err != nil {
			result.Allowed = false
			result.Diagnostics = append(result.Diagnostics, rootlineDiagnostic(err))
			continue
		}
		change.Applied = true
		if _, err := client.ValidateOne(ctx, filepath.Join(cfg.RoadmapRoot, filepath.FromSlash(change.Path))); err != nil {
			result.Allowed = false
			result.Diagnostics = append(result.Diagnostics, rootlineDiagnostic(err))
		}
	}
	postOptions := options
	postOptions.Repo = cfg.RepoRoot
	postOptions.RoadmapRoot = cfg.RoadmapRootRel
	postcheck := runCheck(ctx, postOptions)
	if postcheck.Summary.Errors > 0 {
		result.Allowed = false
	}
	result.Diagnostics = append(result.Diagnostics, postcheck.Diagnostics...)
	return result
}

func validateTransitionTargetStatus(ctx context.Context, cfg *config.Config, options Options, result roadmap.TransitionResult, path string) roadmap.TransitionResult {
	if result.TargetStatus == "" || !result.Allowed {
		return result
	}
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: cfg.RepoRoot, Timeout: options.Timeout})
	describe, err := client.Describe(ctx, cfg.RoadmapRoot, "schema.estado")
	if err != nil {
		result.Allowed = false
		result.Changes = nil
		result.Diagnostics = append(result.Diagnostics, rootlineDiagnostic(err))
		return result
	}
	for _, status := range contextSchemaValues(describe.Decoded, "estado") {
		if status == result.TargetStatus {
			return result
		}
	}
	result.Allowed = false
	result.Changes = nil
	result.Diagnostics = append(result.Diagnostics, diagnostics.Diagnostic{ID: diagnostics.DiagnosticTransitionStatusUnknown, Severity: diagnostics.SeverityError, Message: "target status is not present in effective schema", Path: path, Details: map[string]any{"target": result.TargetStatus}})
	return result
}

func transitionStatusValue(cfg *config.Config, requested string) string {
	switch requested {
	case "pending":
		return cfg.StatusValues.Pending
	case "specified":
		return cfg.StatusValues.Specified
	case "in-progress", "in_progress":
		return cfg.StatusValues.InProgress
	case "completed":
		return cfg.StatusValues.Completed
	case "blocked":
		return cfg.StatusValues.Blocked
	case "obsolete":
		return cfg.StatusValues.Obsolete
	default:
		return requested
	}
}

func normalizeTransitionPath(roadmapRoot string, taskPath string) string {
	if roadmapRoot != "" && filepath.IsAbs(taskPath) {
		if rel, err := filepath.Rel(roadmapRoot, taskPath); err == nil {
			return filepath.ToSlash(filepath.Clean(rel))
		}
	}
	return filepath.ToSlash(filepath.Clean(strings.TrimPrefix(taskPath, "./")))
}

func newTransitionReport(root string, roadmapRoot string, action string, path string, result roadmap.TransitionResult) transitionReport {
	report := diagnostics.NewReport("roadmapctl/transition", root, roadmapRoot, result.Diagnostics)
	return transitionReport{Version: report.Version, Kind: report.Kind, Summary: report.Summary, Root: report.Root, RoadmapRoot: report.RoadmapRoot, Action: action, Path: path, Allowed: result.Allowed, CurrentStatus: result.CurrentStatus, TargetStatus: result.TargetStatus, Role: result.Role, Reasons: result.Reasons, BlockingDependencies: result.BlockingDependencies, Changes: result.Changes, Diagnostics: report.Diagnostics}
}
