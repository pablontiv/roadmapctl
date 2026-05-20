package diagnostics

import (
	"encoding/json"
	"io"
	"path/filepath"

	diag "github.com/pablontiv/picokit/diag"
)

type Summary = diag.Summary
type Diagnostic = diag.Diagnostic
type Severity = diag.Severity

const (
	ExitOK               = diag.ExitOK
	ExitValidation       = diag.ExitValidation
	ExitUsage            = diag.ExitUsage
	ExitEnvironment      = diag.ExitEnvironment
	ExitInternal         = diag.ExitInternal
	SeverityInfo         = diag.SeverityInfo
	SeverityWarning      = diag.SeverityWarning
	SeverityError        = diag.SeverityError
	SummaryStatusOK      = diag.SummaryStatusOK
	SummaryStatusWarning = diag.SummaryStatusWarning
	SummaryStatusError   = diag.SummaryStatusError
)

type Report struct {
	Version     int          `json:"version"`
	Kind        string       `json:"kind"`
	Summary     Summary      `json:"summary"`
	Root        string       `json:"root"`
	RoadmapRoot string       `json:"roadmap_root"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

func toBase(r Report) diag.Report {
	return diag.Report{
		Version: r.Version, Kind: r.Kind, Summary: r.Summary,
		Root: r.Root, Diagnostics: r.Diagnostics}
}
func NewReport(kind, root, roadmapRoot string, diagnostics []Diagnostic) Report {
	base := diag.NewReport(kind, root, diagnostics)
	return Report{Version: base.Version, Kind: base.Kind, Summary: base.Summary,
		Root: base.Root, RoadmapRoot: filepath.ToSlash(roadmapRoot), Diagnostics: base.Diagnostics}
}

func ExitCode(r Report, strict bool) int {
	return diag.ExitCode(toBase(r), strict)
}

func RenderJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	return enc.Encode(r)
}

func RenderText(w io.Writer, r Report) error {
	return diag.RenderText(w, toBase(r))
}
