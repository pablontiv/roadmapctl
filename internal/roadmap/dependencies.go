package roadmap

import (
	"context"
	"errors"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	roadmaplint "github.com/pablontiv/roadmapctl/internal/lint"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
)

const (
	DiagnosticGraphCycle                 = "RMC_GRAPH_CYCLE"
	DiagnosticStatusUnknown              = "RMC_STATUS_UNKNOWN"
	DiagnosticTypeUnknown                = "RMC_STATUS_TYPE_UNKNOWN"
	DiagnosticRootlineValidateFailed     = "RMC_ROOTLINE_VALIDATE_FAILED"
	DiagnosticRootlineDescribeFailed     = "RMC_ROOTLINE_DESCRIBE_FAILED"
	DiagnosticRootlineQueryFailed        = "RMC_ROOTLINE_QUERY_FAILED"
	DiagnosticRootlineGraphFailed        = "RMC_ROOTLINE_GRAPH_FAILED"
	DiagnosticConfigStatusSchemaMismatch = "RMC_CONFIG_STATUS_SCHEMA_MISMATCH"
)

type RootlineClient interface {
	Validate(ctx context.Context, paths ...string) (*rootlinecli.JSONResult, error)
	Describe(ctx context.Context, target string, fields ...string) (*rootlinecli.JSONResult, error)
	Query(ctx context.Context, root string, wheres ...string) (*rootlinecli.JSONResult, error)
	Graph(ctx context.Context, root string, wheres ...string) (*rootlinecli.JSONResult, error)
}

type OperationalStatus struct {
	Source string
	Value  string
	Path   string
}

type RootlineCheckOptions struct {
	RoadmapRoot         string
	LeafFilter          string
	AllowedStatuses     []string
	OperationalStatuses []OperationalStatus
}

func CheckRootline(ctx context.Context, cfg *config.Config, client RootlineClient, options RootlineCheckOptions) ([]Diagnostic, error) {
	var found []Diagnostic

	validateResult, err := client.Validate(ctx, "--all", options.RoadmapRoot)
	if err != nil {
		if validateResult == nil {
			found = append(found, rootlineOperationDiagnostic("validate", err))
		} else {
			parsed := validateDiagnostics(validateResult.Decoded)
			if len(parsed) == 0 {
				found = append(found, rootlineOperationDiagnostic("validate", err))
			} else {
				found = append(found, addRootlineErrorDetails("validate", err, parsed)...)
			}
		}
		if isMissingRootline(err) {
			return found, nil
		}
	} else {
		found = append(found, validateDiagnostics(validateResult.Decoded)...)
	}

	describeResult, err := client.Describe(ctx, ensureDirPath(options.RoadmapRoot))
	schemaStatuses := []string(nil)
	schemaTypes := []string(nil)
	if err != nil {
		found = append(found, rootlineOperationDiagnostic("describe", err))
	} else {
		schemaStatuses = extractStatusValues(describeResult.Decoded, cfg)
		schemaTypes = extractTypeValues(describeResult.Decoded, cfg)
		found = append(found, operationalStatusDiagnostics(options.OperationalStatuses, schemaStatuses)...)
		found = append(found, roadmaplint.CheckOutcomeSchemaCompatibility(describeResult.Decoded)...)
	}

	queryFilter := cfg.Fields.RecordType + ` == "` + cfg.Fields.TaskValue + `"`
	queryResult, err := client.Query(ctx, options.RoadmapRoot, options.LeafFilter, queryFilter)
	if err != nil {
		found = append(found, rootlineOperationDiagnostic("query", err))
	} else {
		found = append(found, statusDiagnostics(queryResult.Decoded, cfg, options.AllowedStatuses, schemaStatuses, schemaTypes)...)
	}

	graphResult, err := client.Graph(ctx, options.RoadmapRoot, options.LeafFilter)
	if err != nil {
		found = append(found, rootlineOperationDiagnostic("graph", err))
	} else {
		found = append(found, graphDiagnostics(cfg, graphResult.Decoded)...)
	}

	return found, nil
}

func validateDiagnostics(decoded map[string]any) []Diagnostic {
	invalid := numberAt(decoded, "summary", "invalid")
	if invalid == 0 {
		invalid = numberAt(decoded, "summary", "invalid_count")
	}
	if invalid == 0 {
		return nil
	}
	return []Diagnostic{{
		ID:       DiagnosticRootlineValidateFailed,
		Severity: diagnostics.SeverityError,
		Message:  "rootline validation reported invalid roadmap records",
		Details:  map[string]any{"invalid": invalid},
	}}
}

func graphDiagnostics(cfg *config.Config, decoded map[string]any) []Diagnostic {
	var found []Diagnostic
	for _, cycle := range arrayValue(decoded["cycles"]) {
		found = append(found, Diagnostic{
			ID:       DiagnosticGraphCycle,
			Severity: diagnostics.SeverityError,
			Message:  "roadmap dependency graph contains a cycle",
			Details:  map[string]any{"cycle": cycle},
		})
	}
	for _, broken := range arrayValue(decoded["broken_links"]) {
		link, ok := broken.(map[string]any)
		if !ok {
			continue
		}
		if stringField(link, "type") != cfg.Fields.DependencyLink {
			continue
		}
		found = append(found, Diagnostic{
			ID:       diagnostics.DiagnosticInvalidBlockedBy,
			Severity: diagnostics.SeverityError,
			Message:  "blocked_by link is broken or invalid",
			Path:     stringField(link, "source"),
			Details: map[string]any{
				"target": stringField(link, "target"),
				"line":   link["line"],
			},
		})
	}
	return found
}

func rootlineOperationDiagnostic(operation string, err error) Diagnostic {
	var rootlineErr *rootlinecli.Error
	if errors.As(err, &rootlineErr) {
		if rootlineErr.Kind == rootlinecli.ErrorMissingBinary {
			return rootlineErr.Diagnostic()
		}
		exitCode := rootlineErr.ExitCode
		if operation == "validate" && rootlineErr.Kind == rootlinecli.ErrorExecution {
			exitCode = diagnostics.ExitValidation
		}
		return Diagnostic{
			ID:       rootlineDiagnosticID(operation),
			Severity: diagnostics.SeverityError,
			Message:  rootlineErr.Message,
			Path:     rootlineErr.Path,
			Details: map[string]any{
				"operation": operation,
				"kind":      string(rootlineErr.Kind),
				"stderr":    rootlineErr.Stderr,
			},
			ExitCode: exitCode,
		}
	}
	return Diagnostic{
		ID:       rootlineDiagnosticID(operation),
		Severity: diagnostics.SeverityError,
		Message:  err.Error(),
		Details:  map[string]any{"operation": operation},
		ExitCode: diagnostics.ExitEnvironment,
	}
}

func addRootlineErrorDetails(operation string, err error, parsed []Diagnostic) []Diagnostic {
	operationDiagnostic := rootlineOperationDiagnostic(operation, err)
	for i := range parsed {
		if parsed[i].Details == nil {
			parsed[i].Details = map[string]any{}
		}
		for key, value := range operationDiagnostic.Details {
			parsed[i].Details[key] = value
		}
		if parsed[i].ExitCode == 0 {
			parsed[i].ExitCode = operationDiagnostic.ExitCode
		}
	}
	return parsed
}

func isMissingRootline(err error) bool {
	var rootlineErr *rootlinecli.Error
	return errors.As(err, &rootlineErr) && rootlineErr.Kind == rootlinecli.ErrorMissingBinary
}

func rootlineDiagnosticID(operation string) string {
	switch operation {
	case "validate":
		return DiagnosticRootlineValidateFailed
	case "describe":
		return DiagnosticRootlineDescribeFailed
	case "query":
		return DiagnosticRootlineQueryFailed
	case "graph":
		return DiagnosticRootlineGraphFailed
	default:
		return "RMC_ROOTLINE_ERROR"
	}
}

func numberAt(decoded map[string]any, keys ...string) int {
	var current any = decoded
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return 0
		}
		current = m[key]
	}
	switch v := current.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}

func ensureDirPath(path string) string {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, `\`) {
		return path
	}
	return path + "/"
}
