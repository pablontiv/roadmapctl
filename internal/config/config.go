package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/fsx"
	"github.com/pelletier/go-toml/v2"
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
	Warnings       []Warning

	DoneStatuses   []string
	ActiveStatuses []string
	StatusValues   StatusValues
	LeafFilter     string

	OutcomeCloseVerify   []string
	PRMergeStrategy      string
	CommitStyle          string
	AutoPush             bool
	RequiredCodeCoverage float64

	LoopMaxTasks           int
	Parallel               bool
	Autonomy               string
	CompactAfterTaskCommit bool
	PRMode                 bool
}

type Warning struct {
	Code    string
	Message string
	Path    string
}

const WarnConfigConflict = "RMC_CONFIG_SOURCE_CONFLICT"

type MigrationPlan struct {
	SourcePath string
	TargetPath string
	Content    string
}

type StatusValues struct {
	Pending    string `json:"pending"`
	Specified  string `json:"specified"`
	InProgress string `json:"in_progress"`
	Completed  string `json:"completed"`
	Blocked    string `json:"blocked"`
	Obsolete   string `json:"obsolete"`
}

func Load(repo string, opts Options) (*Config, error) {
	absRepo, err := filepath.Abs(repo)
	if err != nil {
		return nil, &Error{Code: ErrConfigParse, Message: "resolve repo root", ExitCode: 2, Cause: err}
	}
	absRepo = filepath.Clean(absRepo)
	cfg := defaultConfig(absRepo)

	legacyPath := filepath.Join(absRepo, ".claude", "roadmap.local.md")
	roadmapRoot := opts.RoadmapRoot
	if strings.TrimSpace(roadmapRoot) == "" {
		roadmapRoot = filepath.ToSlash(filepath.Join("docs", "roadmap"))
	}

	tomlRoadmapRoot := roadmapRoot
	tomlPath := filepath.Join(absRepo, filepath.FromSlash(normalizeSeparators(tomlRoadmapRoot)), ".roadmapctl.toml")
	if fileExists(tomlPath) {
		cfg.ConfigPath = tomlPath
		if err := loadTOMLConfig(cfg, tomlPath); err != nil {
			return nil, err
		}
		roadmapRoot = filepath.ToSlash(filepath.Dir(strings.TrimPrefix(tomlPath, absRepo+string(filepath.Separator))))
		if fileExists(legacyPath) {
			if err := os.Remove(legacyPath); err != nil {
				return nil, &Error{Code: ErrConfigParse, Message: "remove legacy roadmap config after TOML load", Path: legacyPath, ExitCode: 2, Cause: err}
			}
		}
	} else if fileExists(legacyPath) {
		fields, err := loadLegacyFields(legacyPath)
		if err != nil {
			return nil, err
		}
		applyFields(cfg, fields)
		if err := validateConfig(cfg, legacyPath); err != nil {
			return nil, err
		}
		if strings.TrimSpace(opts.RoadmapRoot) == "" {
			roadmapRoot = stringValue(fields["roadmap-root"])
		}
		if strings.TrimSpace(roadmapRoot) == "" {
			return nil, &Error{Code: ErrRoadmapRootMissing, Message: "roadmap-root is required", Path: legacyPath, ExitCode: 2}
		}
		absRoadmapRoot, _, err := fsx.ResolveInside(absRepo, roadmapRoot)
		if err != nil {
			return nil, &Error{Code: ErrRoadmapRootEscape, Message: "roadmap-root must resolve inside repo", Path: legacyPath, ExitCode: 2, Cause: err}
		}
		if err := os.MkdirAll(absRoadmapRoot, 0o755); err != nil {
			return nil, &Error{Code: ErrConfigParse, Message: "create roadmap root for TOML migration", Path: absRoadmapRoot, ExitCode: 2, Cause: err}
		}
		migratedPath := filepath.Join(absRoadmapRoot, ".roadmapctl.toml")
		if err := os.WriteFile(migratedPath, []byte(renderTOMLConfig(cfg)), 0o644); err != nil {
			return nil, &Error{Code: ErrConfigParse, Message: "write migrated roadmapctl TOML", Path: migratedPath, ExitCode: 2, Cause: err}
		}
		migratedCfg := defaultConfig(absRepo)
		migratedCfg.ConfigPath = migratedPath
		if err := loadTOMLConfig(migratedCfg, migratedPath); err != nil {
			return nil, err
		}
		if err := os.Remove(legacyPath); err != nil {
			return nil, &Error{Code: ErrConfigParse, Message: "remove legacy roadmap config after TOML migration", Path: legacyPath, ExitCode: 2, Cause: err}
		}
		cfg = migratedCfg
		roadmapRoot = filepath.ToSlash(filepath.Dir(strings.TrimPrefix(migratedPath, absRepo+string(filepath.Separator))))
	} else if !roadmapRootExists(absRepo, roadmapRoot) {
		missingPath := tomlPath
		if strings.TrimSpace(opts.RoadmapRoot) == "" {
			missingPath = legacyPath
		}
		return nil, &Error{Code: ErrConfigMissing, Message: "roadmap config not found", Path: missingPath, ExitCode: 2, Cause: os.ErrNotExist}
	} else {
		cfg.ConfigPath = tomlPath
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

type tomlConfig struct {
	DoneStatuses           []string         `toml:"done_statuses"`
	ActiveStatuses         []string         `toml:"active_statuses"`
	LeafFilter             string           `toml:"leaf_filter"`
	OutcomeCloseVerify     []string         `toml:"outcome_close_verify"`
	PRMergeStrategy        string           `toml:"pr_merge_strategy"`
	CommitStyle            string           `toml:"commit_style"`
	AutoPush               *bool            `toml:"auto_push"`
	RequiredCodeCoverage   *float64         `toml:"required_code_coverage"`
	LoopMaxTasks           *int             `toml:"loop_max_tasks"`
	Parallel               *bool            `toml:"parallel"`
	Autonomy               string           `toml:"autonomy"`
	CompactAfterTaskCommit *bool            `toml:"compact_after_task_commit"`
	PRMode                 *bool            `toml:"pr_mode"`
	StatusValues           tomlStatusValues `toml:"status_values"`
}

type tomlStatusValues struct {
	Pending    string `toml:"pending"`
	Specified  string `toml:"specified"`
	InProgress string `toml:"in_progress"`
	Completed  string `toml:"completed"`
	Blocked    string `toml:"blocked"`
	Obsolete   string `toml:"obsolete"`
}

func LegacyMigrationPlan(repo string, opts Options) (MigrationPlan, error) {
	absRepo, err := filepath.Abs(repo)
	if err != nil {
		return MigrationPlan{}, &Error{Code: ErrConfigParse, Message: "resolve repo root", ExitCode: 2, Cause: err}
	}
	absRepo = filepath.Clean(absRepo)
	legacyPath := filepath.Join(absRepo, ".claude", "roadmap.local.md")
	fields, err := loadLegacyFields(legacyPath)
	if err != nil {
		return MigrationPlan{}, err
	}
	cfg := defaultConfig(absRepo)
	applyFields(cfg, fields)
	roadmapRoot := stringValue(fields["roadmap-root"])
	if opts.RoadmapRoot != "" {
		roadmapRoot = opts.RoadmapRoot
	}
	absRoadmapRoot, _, err := fsx.ResolveInside(absRepo, roadmapRoot)
	if err != nil {
		return MigrationPlan{}, &Error{Code: ErrRoadmapRootEscape, Message: "roadmap-root must resolve inside repo", Path: legacyPath, ExitCode: 2, Cause: err}
	}
	return MigrationPlan{
		SourcePath: legacyPath,
		TargetPath: filepath.Join(absRoadmapRoot, ".roadmapctl.toml"),
		Content:    renderTOMLConfig(cfg),
	}, nil
}

func loadTOMLConfig(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return &Error{Code: ErrConfigParse, Message: "read roadmapctl config", Path: path, ExitCode: 2, Cause: err}
	}
	var decoded tomlConfig
	if err := toml.Unmarshal(data, &decoded); err != nil {
		return &Error{Code: ErrConfigParse, Message: "parse roadmapctl TOML: " + err.Error(), Path: path, ExitCode: 2, Cause: err}
	}
	applyTOMLConfig(cfg, decoded)
	if err := validateConfig(cfg, path); err != nil {
		return err
	}
	return nil
}

func loadLegacyFields(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &Error{Code: ErrConfigMissing, Message: "roadmap config not found", Path: path, ExitCode: 2, Cause: err}
		}
		return nil, &Error{Code: ErrConfigParse, Message: "read roadmap config", Path: path, ExitCode: 2, Cause: err}
	}
	fields, err := parseFrontmatter(data)
	if err != nil {
		return nil, &Error{Code: ErrConfigParse, Message: err.Error(), Path: path, ExitCode: 2, Cause: err}
	}
	return fields, nil
}

func applyTOMLConfig(cfg *Config, decoded tomlConfig) {
	if decoded.DoneStatuses != nil {
		cfg.DoneStatuses = append([]string(nil), decoded.DoneStatuses...)
	}
	if decoded.ActiveStatuses != nil {
		cfg.ActiveStatuses = append([]string(nil), decoded.ActiveStatuses...)
	}
	if decoded.LeafFilter != "" {
		cfg.LeafFilter = decoded.LeafFilter
	}
	if decoded.OutcomeCloseVerify != nil {
		cfg.OutcomeCloseVerify = append([]string(nil), decoded.OutcomeCloseVerify...)
	}
	if decoded.PRMergeStrategy != "" {
		cfg.PRMergeStrategy = decoded.PRMergeStrategy
	}
	if decoded.CommitStyle != "" {
		cfg.CommitStyle = decoded.CommitStyle
	}
	if decoded.AutoPush != nil {
		cfg.AutoPush = *decoded.AutoPush
	}
	if decoded.RequiredCodeCoverage != nil {
		cfg.RequiredCodeCoverage = *decoded.RequiredCodeCoverage
	}
	if decoded.LoopMaxTasks != nil {
		cfg.LoopMaxTasks = *decoded.LoopMaxTasks
	}
	if decoded.Parallel != nil {
		cfg.Parallel = *decoded.Parallel
	}
	if decoded.Autonomy != "" {
		cfg.Autonomy = decoded.Autonomy
	}
	if decoded.CompactAfterTaskCommit != nil {
		cfg.CompactAfterTaskCommit = *decoded.CompactAfterTaskCommit
	}
	if decoded.PRMode != nil {
		cfg.PRMode = *decoded.PRMode
	}
	if decoded.StatusValues.Pending != "" {
		cfg.StatusValues.Pending = decoded.StatusValues.Pending
	}
	if decoded.StatusValues.Specified != "" {
		cfg.StatusValues.Specified = decoded.StatusValues.Specified
	}
	if decoded.StatusValues.InProgress != "" {
		cfg.StatusValues.InProgress = decoded.StatusValues.InProgress
	}
	if decoded.StatusValues.Completed != "" {
		cfg.StatusValues.Completed = decoded.StatusValues.Completed
	}
	if decoded.StatusValues.Blocked != "" {
		cfg.StatusValues.Blocked = decoded.StatusValues.Blocked
	}
	if decoded.StatusValues.Obsolete != "" {
		cfg.StatusValues.Obsolete = decoded.StatusValues.Obsolete
	}
}

func renderTOMLConfig(cfg *Config) string {
	var b strings.Builder
	writeStringList(&b, "done_statuses", cfg.DoneStatuses)
	writeStringList(&b, "active_statuses", cfg.ActiveStatuses)
	fmt.Fprintf(&b, "leaf_filter = '%s'\n", cfg.LeafFilter)
	writeStringList(&b, "outcome_close_verify", cfg.OutcomeCloseVerify)
	fmt.Fprintf(&b, "pr_merge_strategy = '%s'\n", cfg.PRMergeStrategy)
	fmt.Fprintf(&b, "commit_style = '%s'\n", cfg.CommitStyle)
	fmt.Fprintf(&b, "auto_push = %t\n", cfg.AutoPush)
	fmt.Fprintf(&b, "required_code_coverage = %.1f\n", cfg.RequiredCodeCoverage)
	fmt.Fprintf(&b, "loop_max_tasks = %d\n", cfg.LoopMaxTasks)
	fmt.Fprintf(&b, "parallel = %t\n", cfg.Parallel)
	fmt.Fprintf(&b, "autonomy = '%s'\n", cfg.Autonomy)
	fmt.Fprintf(&b, "compact_after_task_commit = %t\n", cfg.CompactAfterTaskCommit)
	fmt.Fprintf(&b, "pr_mode = %t\n\n", cfg.PRMode)
	b.WriteString("[status_values]\n")
	fmt.Fprintf(&b, "pending = '%s'\n", cfg.StatusValues.Pending)
	fmt.Fprintf(&b, "specified = '%s'\n", cfg.StatusValues.Specified)
	fmt.Fprintf(&b, "in_progress = '%s'\n", cfg.StatusValues.InProgress)
	fmt.Fprintf(&b, "completed = '%s'\n", cfg.StatusValues.Completed)
	fmt.Fprintf(&b, "blocked = '%s'\n", cfg.StatusValues.Blocked)
	fmt.Fprintf(&b, "obsolete = '%s'\n", cfg.StatusValues.Obsolete)
	return b.String()
}

func writeStringList(b *strings.Builder, key string, values []string) {
	fmt.Fprintf(b, "%s = [", key)
	for i, value := range values {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(b, "'%s'", value)
	}
	b.WriteString("]\n")
}

func configDiffers(left *Config, right *Config) bool {
	return !stringSlicesEqual(left.DoneStatuses, right.DoneStatuses) ||
		!stringSlicesEqual(left.ActiveStatuses, right.ActiveStatuses) ||
		left.LeafFilter != right.LeafFilter ||
		!stringSlicesEqual(left.OutcomeCloseVerify, right.OutcomeCloseVerify) ||
		left.PRMergeStrategy != right.PRMergeStrategy ||
		left.CommitStyle != right.CommitStyle ||
		left.AutoPush != right.AutoPush ||
		left.RequiredCodeCoverage != right.RequiredCodeCoverage ||
		left.LoopMaxTasks != right.LoopMaxTasks ||
		left.Parallel != right.Parallel ||
		left.Autonomy != right.Autonomy ||
		left.CompactAfterTaskCommit != right.CompactAfterTaskCommit ||
		left.PRMode != right.PRMode ||
		left.StatusValues != right.StatusValues
}

func stringSlicesEqual(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func roadmapRootExists(repo string, roadmapRoot string) bool {
	root, _, err := fsx.ResolveInside(repo, roadmapRoot)
	if err != nil {
		return false
	}
	info, err := os.Stat(root)
	return err == nil && info.IsDir()
}

func normalizeSeparators(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
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
		LeafFilter:             "isIndex == false",
		OutcomeCloseVerify:     []string{},
		PRMergeStrategy:        "squash",
		CommitStyle:            "conventional",
		AutoPush:               true,
		RequiredCodeCoverage:   85.0,
		LoopMaxTasks:           0,
		Parallel:               true,
		Autonomy:               "until_done",
		CompactAfterTaskCommit: true,
		PRMode:                 false,
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
	if v, ok := floatValue(fields["required-code-coverage"]); ok {
		cfg.RequiredCodeCoverage = v
	}
	if v, ok := intValue(fields["loop-max-tasks"]); ok {
		cfg.LoopMaxTasks = v
	}
	if v, ok := boolValue(fields["parallel"]); ok {
		cfg.Parallel = v
	}
	if v, ok := stringValueOK(fields["autonomy"]); ok {
		cfg.Autonomy = v
	}
	if v, ok := boolValue(fields["compact-after-task-commit"]); ok {
		cfg.CompactAfterTaskCommit = v
	}
	if v, ok := boolValue(fields["pr-mode"]); ok {
		cfg.PRMode = v
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
	var number float64
	if _, err := fmt.Sscanf(value, "%f", &number); err == nil {
		return number
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

func intValue(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}

func floatValue(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func validateConfig(cfg *Config, path string) error {
	if cfg.RequiredCodeCoverage < 0 || cfg.RequiredCodeCoverage > 100 {
		return &Error{Code: ErrConfigParse, Message: "required_code_coverage must be between 0 and 100", Path: path, ExitCode: 2}
	}
	if cfg.LoopMaxTasks < 0 {
		return &Error{Code: ErrConfigParse, Message: "loop_max_tasks must be greater than or equal to 0", Path: path, ExitCode: 2}
	}
	switch cfg.Autonomy {
	case "manual", "supervised", "until_done":
		return nil
	default:
		return &Error{Code: ErrConfigParse, Message: "autonomy must be one of manual, supervised, until_done", Path: path, ExitCode: 2}
	}
}
