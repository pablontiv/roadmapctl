package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/fsx"
)

const (
	ErrConfigMissing      = "RMC_CONFIG_MISSING"
	ErrConfigParse        = "RMC_CONFIG_PARSE"
	ErrRoadmapRootMissing = "RMC_CONFIG_ROADMAP_ROOT_MISSING"
	ErrRoadmapRootEscape  = "RMC_CONFIG_ROADMAP_ROOT_ESCAPE"
)

type Error struct {
	Code     string
	Message  string
	Path     string
	ExitCode int
	Cause    error
}

func (e *Error) Error() string {
	if e.Path == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s: %s", e.Code, e.Path, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

type Options struct {
	RoadmapRoot string
}

type Config struct {
	RepoRoot       string
	ConfigPath     string
	RoadmapRoot    string
	RoadmapRootRel string

	DoneStatuses   []string
	ActiveStatuses []string
	StatusValues   StatusValues
	LeafFilter     string

	OutcomeCloseVerify []string
	PRMergeStrategy    string
	CommitStyle        string
	AutoPush           bool
}

type StatusValues struct {
	Pending    string
	Specified  string
	InProgress string
	Completed  string
	Blocked    string
	Obsolete   string
}

func Load(repo string, opts Options) (*Config, error) {
	absRepo, err := filepath.Abs(repo)
	if err != nil {
		return nil, &Error{Code: ErrConfigParse, Message: "resolve repo root", ExitCode: 2, Cause: err}
	}
	absRepo = filepath.Clean(absRepo)
	cfg := defaultConfig(absRepo)
	cfg.ConfigPath = filepath.Join(absRepo, ".claude", "roadmap.local.md")

	data, err := os.ReadFile(cfg.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &Error{Code: ErrConfigMissing, Message: "roadmap config not found", Path: cfg.ConfigPath, ExitCode: 2, Cause: err}
		}
		return nil, &Error{Code: ErrConfigParse, Message: "read roadmap config", Path: cfg.ConfigPath, ExitCode: 2, Cause: err}
	}

	fields, err := parseFrontmatter(data)
	if err != nil {
		return nil, &Error{Code: ErrConfigParse, Message: err.Error(), Path: cfg.ConfigPath, ExitCode: 2, Cause: err}
	}
	applyFields(cfg, fields)

	roadmapRoot := stringValue(fields["roadmap-root"])
	if opts.RoadmapRoot != "" {
		roadmapRoot = opts.RoadmapRoot
	}
	if strings.TrimSpace(roadmapRoot) == "" {
		return nil, &Error{Code: ErrRoadmapRootMissing, Message: "roadmap-root is required", Path: cfg.ConfigPath, ExitCode: 2}
	}

	absRoadmapRoot, relRoadmapRoot, err := fsx.ResolveInside(absRepo, roadmapRoot)
	if err != nil {
		return nil, &Error{Code: ErrRoadmapRootEscape, Message: "roadmap-root must resolve inside repo", Path: cfg.ConfigPath, ExitCode: 2, Cause: err}
	}
	cfg.RoadmapRoot = absRoadmapRoot
	cfg.RoadmapRootRel = relRoadmapRoot

	return cfg, nil
}

func defaultConfig(repo string) *Config {
	return &Config{
		RepoRoot:       repo,
		DoneStatuses:   []string{"Completed", "Obsolete"},
		ActiveStatuses: []string{"Pending", "Specified", "In Progress"},
		StatusValues: StatusValues{
			Pending:    "Pending",
			Specified:  "Specified",
			InProgress: "In Progress",
			Completed:  "Completed",
			Blocked:    "Blocked",
			Obsolete:   "Obsolete",
		},
		LeafFilter:         "isIndex == false",
		OutcomeCloseVerify: []string{},
		PRMergeStrategy:    "squash",
		CommitStyle:        "conventional",
		AutoPush:           true,
	}
}

func applyFields(cfg *Config, fields map[string]any) {
	if v, ok := stringSliceValue(fields["done-statuses"]); ok {
		cfg.DoneStatuses = v
	}
	if v, ok := stringSliceValue(fields["active-statuses"]); ok {
		cfg.ActiveStatuses = v
	}
	if v, ok := stringValueOK(fields["leaf-filter"]); ok {
		cfg.LeafFilter = v
	}
	if v, ok := stringSliceValue(fields["outcome-close-verify"]); ok {
		cfg.OutcomeCloseVerify = v
	}
	if v, ok := stringValueOK(fields["pr-merge-strategy"]); ok {
		cfg.PRMergeStrategy = v
	}
	if v, ok := stringValueOK(fields["commit-style"]); ok {
		cfg.CommitStyle = v
	}
	if v, ok := boolValue(fields["auto-push"]); ok {
		cfg.AutoPush = v
	}
	if values, ok := fields["status-values"].(map[string]any); ok {
		if v, ok := stringValueOK(values["pending"]); ok {
			cfg.StatusValues.Pending = v
		}
		if v, ok := stringValueOK(values["specified"]); ok {
			cfg.StatusValues.Specified = v
		}
		if v, ok := stringValueOK(values["in-progress"]); ok {
			cfg.StatusValues.InProgress = v
		}
		if v, ok := stringValueOK(values["completed"]); ok {
			cfg.StatusValues.Completed = v
		}
		if v, ok := stringValueOK(values["blocked"]); ok {
			cfg.StatusValues.Blocked = v
		}
		if v, ok := stringValueOK(values["obsolete"]); ok {
			cfg.StatusValues.Obsolete = v
		}
	}
}

func parseFrontmatter(data []byte) (map[string]any, error) {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil, fmt.Errorf("missing YAML frontmatter")
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		return nil, fmt.Errorf("unterminated YAML frontmatter")
	}

	return parseYAMLLines(lines[1:end])
}

func parseYAMLLines(lines []string) (map[string]any, error) {
	result := map[string]any{}
	var currentMap string

	for _, raw := range lines {
		if strings.TrimSpace(raw) == "" || strings.HasPrefix(strings.TrimSpace(raw), "#") {
			continue
		}
		indent := len(raw) - len(strings.TrimLeft(raw, " "))
		line := strings.TrimSpace(raw)
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid YAML line %q", raw)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if indent == 0 {
			currentMap = ""
			if value == "" {
				child := map[string]any{}
				result[key] = child
				currentMap = key
				continue
			}
			result[key] = parseScalar(value)
			continue
		}

		if currentMap == "" {
			return nil, fmt.Errorf("nested value without parent: %q", raw)
		}
		child, ok := result[currentMap].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("parent %q is not a map", currentMap)
		}
		child[key] = parseScalar(value)
	}
	return result, nil
}

func parseScalar(value string) any {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return parseInlineStringList(value)
	}
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	return unquote(value)
}

func parseInlineStringList(value string) []string {
	inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(value, "["), "]"))
	if inner == "" {
		return []string{}
	}
	parts := strings.Split(inner, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		items = append(items, unquote(strings.TrimSpace(part)))
	}
	return items
}

func unquote(value string) string {
	if len(value) >= 2 {
		first := value[0]
		last := value[len(value)-1]
		if (first == '\'' && last == '\'') || (first == '"' && last == '"') {
			return value[1 : len(value)-1]
		}
	}
	return value
}

func stringValue(value any) string {
	v, _ := stringValueOK(value)
	return v
}

func stringValueOK(value any) (string, bool) {
	v, ok := value.(string)
	return v, ok
}

func stringSliceValue(value any) ([]string, bool) {
	v, ok := value.([]string)
	if !ok {
		return nil, false
	}
	copyValue := append([]string(nil), v...)
	return copyValue, true
}

func boolValue(value any) (bool, bool) {
	v, ok := value.(bool)
	return v, ok
}
