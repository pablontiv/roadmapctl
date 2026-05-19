package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	roadmaplint "github.com/pablontiv/roadmapctl/internal/lint"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
)

const (
	diagnosticRoadmapRootMissing = "RMC_CONFIG_ROADMAP_ROOT_MISSING"
	diagnosticStemMissing        = "RMC_CONFIG_STEM_MISSING"
	diagnosticDoctorPath         = "RMC_ENV_PATH"
)

func runDoctor(ctx context.Context, options Options) diagnostics.Report {
	repoRoot := absoluteClean(options.Repo)
	var roadmapRoot string
	var found []diagnostics.Diagnostic

	cfg, err := config.Load(options.Repo)
	if err != nil {
		found = append(found, configDiagnostic(repoRoot, err))
		return diagnostics.NewReport("roadmapctl/doctor", repoRoot, roadmapRoot, found)
	}

	found = append(found, configWarnings(cfg)...)
	found = append(found, inspectRoadmapPaths(cfg)...)

	client := rootlinecli.New(rootlinecli.Options{
		Binary:  options.Rootline,
		Dir:     cfg.RepoRoot,
		Timeout: options.Timeout,
	})
	version, err := client.Version(ctx)
	if err != nil {
		found = append(found, rootlineDiagnostic(err))
	} else {
		found = append(found, diagnostics.Diagnostic{
			ID:       diagnosticDoctorPath,
			Severity: diagnostics.SeverityInfo,
			Message:  "roadmapctl doctor paths resolved",
			Path:     cfg.RoadmapRootRel,
			Details: map[string]any{
				"config":           relToRoot(cfg.RepoRoot, cfg.ConfigPath),
				"roadmap_root":     cfg.RoadmapRootRel,
				"rootline_version": strings.TrimSpace(string(version.Stdout)),
			},
		})
		describe, err := client.Describe(ctx, ensureRootlineDirPath(cfg.RoadmapRoot))
		if err != nil {
			found = append(found, rootlineDiagnostic(err))
		} else {
			found = append(found, roadmaplint.CheckOutcomeSchemaCompatibility(describe.Decoded)...)
		}
	}

	return diagnostics.NewReport("roadmapctl/doctor", cfg.RepoRoot, cfg.RoadmapRoot, found)
}

func inspectRoadmapPaths(cfg *config.Config) []diagnostics.Diagnostic {
	var found []diagnostics.Diagnostic
	if info, err := os.Stat(cfg.RoadmapRoot); err != nil {
		if os.IsNotExist(err) {
			found = append(found, diagnostics.Diagnostic{
				ID:       diagnosticRoadmapRootMissing,
				Severity: diagnostics.SeverityError,
				Message:  "roadmap root not found",
				Path:     cfg.RoadmapRootRel,
				ExitCode: diagnostics.ExitUsage,
			})
		} else {
			found = append(found, diagnostics.Diagnostic{
				ID:       diagnosticRoadmapRootMissing,
				Severity: diagnostics.SeverityError,
				Message:  fmt.Sprintf("inspect roadmap root: %v", err),
				Path:     cfg.RoadmapRootRel,
				ExitCode: diagnostics.ExitUsage,
			})
		}
	} else if !info.IsDir() {
		found = append(found, diagnostics.Diagnostic{
			ID:       diagnosticRoadmapRootMissing,
			Severity: diagnostics.SeverityError,
			Message:  "roadmap root is not a directory",
			Path:     cfg.RoadmapRootRel,
			ExitCode: diagnostics.ExitUsage,
		})
	}

	stemPath := filepath.Join(cfg.RoadmapRoot, ".stem")
	if info, err := os.Stat(stemPath); err != nil {
		if os.IsNotExist(err) {
			found = append(found, diagnostics.Diagnostic{
				ID:       diagnosticStemMissing,
				Severity: diagnostics.SeverityError,
				Message:  "roadmap root is missing .stem",
				Path:     relToRoot(cfg.RepoRoot, stemPath),
				ExitCode: diagnostics.ExitUsage,
			})
		} else {
			found = append(found, diagnostics.Diagnostic{
				ID:       diagnosticStemMissing,
				Severity: diagnostics.SeverityError,
				Message:  fmt.Sprintf("inspect roadmap .stem: %v", err),
				Path:     relToRoot(cfg.RepoRoot, stemPath),
				ExitCode: diagnostics.ExitUsage,
			})
		}
	} else if info.IsDir() {
		found = append(found, diagnostics.Diagnostic{
			ID:       diagnosticStemMissing,
			Severity: diagnostics.SeverityError,
			Message:  "roadmap .stem is not a file",
			Path:     relToRoot(cfg.RepoRoot, stemPath),
			ExitCode: diagnostics.ExitUsage,
		})
	}
	return found
}

func configWarnings(cfg *config.Config) []diagnostics.Diagnostic {
	found := make([]diagnostics.Diagnostic, 0, len(cfg.Warnings))
	for _, warning := range cfg.Warnings {
		found = append(found, diagnostics.Diagnostic{
			ID:       warning.Code,
			Severity: diagnostics.SeverityWarning,
			Message:  warning.Message,
			Path:     relToRoot(cfg.RepoRoot, warning.Path),
		})
	}
	return found
}

func configDiagnostic(repoRoot string, err error) diagnostics.Diagnostic {
	var cfgErr *config.Error
	if errors.As(err, &cfgErr) {
		return diagnostics.Diagnostic{
			ID:       cfgErr.Code,
			Severity: diagnostics.SeverityError,
			Message:  cfgErr.Message,
			Path:     relToRoot(repoRoot, cfgErr.Path),
			ExitCode: cfgErr.ExitCode,
		}
	}
	return diagnostics.Diagnostic{
		ID:       "RMC_CONFIG_ERROR",
		Severity: diagnostics.SeverityError,
		Message:  err.Error(),
		ExitCode: diagnostics.ExitUsage,
	}
}

func rootlineDiagnostic(err error) diagnostics.Diagnostic {
	var rootlineErr *rootlinecli.Error
	if errors.As(err, &rootlineErr) {
		return rootlineErr.Diagnostic()
	}
	return diagnostics.Diagnostic{
		ID:       "RMC_ENV_ROOTLINE_ERROR",
		Severity: diagnostics.SeverityError,
		Message:  err.Error(),
		ExitCode: diagnostics.ExitEnvironment,
	}
}

func absoluteClean(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.Clean(abs)
}

func relToRoot(root string, path string) string {
	if path == "" {
		return ""
	}
	rel, err := filepath.Rel(root, path)
	if err != nil || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." || filepath.IsAbs(rel) {
		return filepath.ToSlash(filepath.Clean(path))
	}
	return filepath.ToSlash(rel)
}

func ensureRootlineDirPath(path string) string {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, `\`) {
		return path
	}
	return path + string(filepath.Separator)
}
