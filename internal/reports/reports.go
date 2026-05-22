package reports

import (
	"encoding/json"
	"io"
	"path/filepath"

	diag "github.com/pablontiv/picokit/diag"
)

type Report struct {
	Version     int               `json:"version"`
	Kind        string            `json:"kind"`
	Summary     diag.Summary      `json:"summary"`
	Root        string            `json:"root"`
	RoadmapRoot string            `json:"roadmap_root"`
	Diagnostics []diag.Diagnostic `json:"diagnostics"`
}

func toBase(r Report) diag.Report {
	return diag.Report{
		Version: r.Version, Kind: r.Kind, Summary: r.Summary,
		Root: r.Root, Diagnostics: r.Diagnostics}
}

func NewReport(kind, root, roadmapRoot string, diagnostics []diag.Diagnostic) Report {
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
