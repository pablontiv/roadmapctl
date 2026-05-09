package diagnostics

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	ExitOK          = 0
	ExitValidation  = 1
	ExitUsage       = 2
	ExitEnvironment = 3
	ExitInternal    = 4
)

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

const (
	SummaryStatusOK      = "ok"
	SummaryStatusWarning = "warning"
	SummaryStatusError   = "error"
)

const (
	DiagnosticSingleFileFallback = "RMC_STRUCTURE_SINGLE_FILE_FALLBACK"
	DiagnosticRootlineMissing    = "RMC_ENV_ROOTLINE_MISSING"
	DiagnosticInvalidBlockedBy   = "RMC_GRAPH_INVALID_BLOCKED_BY"
	DiagnosticConfigMissing      = "RMC_CONFIG_MISSING"
)

const (
	DiagnosticLintTaskTableMissing          = "RMC_LINT_TASK_TABLE_MISSING"
	DiagnosticLintTaskTableMissingRow       = "RMC_LINT_TASK_TABLE_MISSING_ROW"
	DiagnosticLintTaskTableStaleRow         = "RMC_LINT_TASK_TABLE_STALE_ROW"
	DiagnosticLintTaskTableInvalidLink      = "RMC_LINT_TASK_TABLE_INVALID_LINK"
	DiagnosticLintTaskSectionMissing        = "RMC_LINT_TASK_SECTION_MISSING"
	DiagnosticLintAcceptanceCriteriaMissing = "RMC_LINT_ACCEPTANCE_CRITERIA_MISSING"
	DiagnosticLintSourceOfTruthEmpty        = "RMC_LINT_SOURCE_OF_TRUTH_EMPTY"
	DiagnosticLintFilenameCaseCollision     = "RMC_LINT_FILENAME_CASE_COLLISION"
	DiagnosticLintFilenameReserved          = "RMC_LINT_FILENAME_RESERVED"
	DiagnosticLintSchemaFieldMissing        = "RMC_LINT_SCHEMA_FIELD_MISSING"
	DiagnosticLintSchemaLinkMissing         = "RMC_LINT_SCHEMA_LINK_MISSING"
)

type Severity string

type Report struct {
	Version     int          `json:"version"`
	Kind        string       `json:"kind"`
	Summary     Summary      `json:"summary"`
	Root        string       `json:"root"`
	RoadmapRoot string       `json:"roadmap_root"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Summary struct {
	Status   string `json:"status"`
	Errors   int    `json:"errors"`
	Warnings int    `json:"warnings"`
	Infos    int    `json:"infos"`
}

type Diagnostic struct {
	ID       string         `json:"id"`
	Severity Severity       `json:"severity"`
	Message  string         `json:"message"`
	Path     string         `json:"path,omitempty"`
	Details  map[string]any `json:"details,omitempty"`
	ExitCode int            `json:"-"`
}

func NewReport(kind string, root string, roadmapRoot string, diagnostics []Diagnostic) Report {
	copied := make([]Diagnostic, len(diagnostics))
	copy(copied, diagnostics)
	return Report{
		Version:     1,
		Kind:        kind,
		Summary:     summarize(copied),
		Root:        root,
		RoadmapRoot: roadmapRoot,
		Diagnostics: copied,
	}
}

func RenderJSON(w io.Writer, report Report) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(report)
}

func RenderText(w io.Writer, report Report) error {
	if _, err := fmt.Fprintf(w, "%s\nstatus: %s\nerrors: %d\nwarnings: %d\ninfos: %d\n", report.Kind, report.Summary.Status, report.Summary.Errors, report.Summary.Warnings, report.Summary.Infos); err != nil {
		return err
	}
	for _, diagnostic := range report.Diagnostics {
		if diagnostic.Path == "" {
			if _, err := fmt.Fprintf(w, "[%s] %s: %s\n", diagnostic.Severity, diagnostic.ID, diagnostic.Message); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprintf(w, "[%s] %s %s: %s\n", diagnostic.Severity, diagnostic.ID, diagnostic.Path, diagnostic.Message); err != nil {
			return err
		}
	}
	return nil
}

func ExitCode(report Report, strict bool) int {
	code := ExitOK
	for _, diagnostic := range report.Diagnostics {
		if diagnostic.Severity == SeverityWarning && strict && code < ExitValidation {
			code = ExitValidation
		}
		if diagnostic.Severity != SeverityError {
			continue
		}
		diagnosticCode := diagnostic.ExitCode
		if diagnosticCode == 0 {
			diagnosticCode = ExitValidation
		}
		if diagnosticCode > code {
			code = diagnosticCode
		}
	}
	return code
}

func summarize(diagnostics []Diagnostic) Summary {
	summary := Summary{Status: SummaryStatusOK}
	for _, diagnostic := range diagnostics {
		switch diagnostic.Severity {
		case SeverityError:
			summary.Errors++
		case SeverityWarning:
			summary.Warnings++
		case SeverityInfo:
			summary.Infos++
		}
	}
	if summary.Errors > 0 {
		summary.Status = SummaryStatusError
	} else if summary.Warnings > 0 {
		summary.Status = SummaryStatusWarning
	}
	return summary
}
