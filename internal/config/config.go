package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pablontiv/picokit/pathsec"
	"github.com/pelletier/go-toml/v2"
)

const (
	ErrConfigMissing      = "RMC_CONFIG_MISSING"
	ErrConfigParse        = "RMC_CONFIG_PARSE"
	ErrRoadmapRootMissing = "RMC_CONFIG_ROADMAP_ROOT_MISSING"
	ErrRoadmapRootEscape  = "RMC_CONFIG_ROADMAP_ROOT_ESCAPE"

	DefaultDependencyLink = "blocked_by"

	roadmapRootDir = "docs/roadmap"
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

type Config struct {
	RepoRoot       string
	ConfigPath     string
	RoadmapRoot    string
	RoadmapRootRel string
	Warnings       []Warning

	DoneStatuses   []string
	ActiveStatuses []string
	StatusValues   StatusValues
	Fields         FieldsConfig
	LeafFilter     string

	OutcomeCloseVerify   []string
	RequiredCodeCoverage float64
	Gitflow              GitflowConfig

	LoopMaxTasks           int
	Parallel               bool
	Autonomy               string
	CompactAfterTaskCommit bool
}

type Warning struct {
	Code    string
	Message string
	Path    string
}

type StatusValues struct {
	Pending    string `json:"pending"`
	Specified  string `json:"specified"`
	InProgress string `json:"in_progress"`
	Completed  string `json:"completed"`
	Blocked    string `json:"blocked"`
	Obsolete   string `json:"obsolete"`
}

type FieldsConfig struct {
	Lifecycle      string `json:"lifecycle"`
	RecordType     string `json:"record_type"`
	TaskValue      string `json:"task_value"`
	OutcomeValue   string `json:"outcome_value"`
	DisplayName    string `json:"display_name"`
	DependencyLink string `json:"dependency_link"`
}

type GitflowConfig struct {
	BaseBranch   string `json:"base_branch"`
	BranchMode   string `json:"branch_mode"`
	BranchStyle  string `json:"branch_style"`
	PrCreate     string `json:"pr_create"`
	PrTitleStyle string `json:"pr_title_style"`
	PrBodyStyle  string `json:"pr_body_style"`
	CommitStyle  string `json:"commit_style"`
	AutoPush     bool   `json:"auto_push"`
}

func Load(repo string) (*Config, error) {
	absRepo, err := filepath.Abs(repo)
	if err != nil {
		return nil, &Error{Code: ErrConfigParse, Message: "resolve repo root", ExitCode: 2, Cause: err}
	}
	absRepo = filepath.Clean(absRepo)
	cfg := defaultConfig(absRepo)

	roadmapRoot := roadmapRootDir
	tomlPath := filepath.Join(absRepo, filepath.FromSlash(roadmapRoot), ".roadmapctl.toml")
	switch {
	case fileExists(tomlPath):
		cfg.ConfigPath = tomlPath
		if err := loadTOMLConfig(cfg, tomlPath); err != nil {
			return nil, err
		}
		roadmapRoot = filepath.ToSlash(filepath.Dir(strings.TrimPrefix(tomlPath, absRepo+string(filepath.Separator))))
	case !roadmapRootExists(absRepo, roadmapRoot):
		return nil, &Error{Code: ErrConfigMissing, Message: "roadmap config not found", Path: tomlPath, ExitCode: 2, Cause: os.ErrNotExist}
	default:
		cfg.ConfigPath = tomlPath
	}

	if strings.TrimSpace(roadmapRoot) == "" {
		return nil, &Error{Code: ErrRoadmapRootMissing, Message: "roadmap-root is required", Path: cfg.ConfigPath, ExitCode: 2}
	}

	absRoadmapRoot, relRoadmapRoot, err := pathsec.ResolveInside(absRepo, roadmapRoot)
	if err != nil {
		return nil, &Error{Code: ErrRoadmapRootEscape, Message: "roadmap-root must resolve inside repo", Path: cfg.ConfigPath, ExitCode: 2, Cause: err}
	}
	cfg.RoadmapRoot = absRoadmapRoot
	cfg.RoadmapRootRel = relRoadmapRoot

	return cfg, nil
}

type tomlConfig struct {
	DoneStatuses           []string          `toml:"done_statuses"`
	ActiveStatuses         []string          `toml:"active_statuses"`
	LeafFilter             string            `toml:"leaf_filter"`
	OutcomeCloseVerify     []string          `toml:"outcome_close_verify"`
	CommitStyle            string            `toml:"commit_style"`
	AutoPush               *bool             `toml:"auto_push"`
	RequiredCodeCoverage   *float64          `toml:"required_code_coverage"`
	LoopMaxTasks           *int              `toml:"loop_max_tasks"`
	Parallel               *bool             `toml:"parallel"`
	Autonomy               string            `toml:"autonomy"`
	CompactAfterTaskCommit *bool             `toml:"compact_after_task_commit"`
	PRMode                 *bool             `toml:"pr_mode"`
	StatusValues           tomlStatusValues  `toml:"status_values"`
	Fields                 tomlFieldsConfig  `toml:"fields"`
	Gitflow                tomlGitflowConfig `toml:"gitflow"`
}

type tomlFieldsConfig struct {
	Lifecycle      string `toml:"lifecycle"`
	RecordType     string `toml:"record_type"`
	TaskValue      string `toml:"task_value"`
	OutcomeValue   string `toml:"outcome_value"`
	DisplayName    string `toml:"display_name"`
	DependencyLink string `toml:"dependency_link"`
}

type tomlStatusValues struct {
	Pending    string `toml:"pending"`
	Specified  string `toml:"specified"`
	InProgress string `toml:"in_progress"`
	Completed  string `toml:"completed"`
	Blocked    string `toml:"blocked"`
	Obsolete   string `toml:"obsolete"`
}

type tomlGitflowConfig struct {
	BaseBranch   string `toml:"base_branch"`
	BranchMode   string `toml:"branch_mode"`
	BranchStyle  string `toml:"branch_style"`
	PrCreate     string `toml:"pr_create"`
	PrTitleStyle string `toml:"pr_title_style"`
	PrBodyStyle  string `toml:"pr_body_style"`
	CommitStyle  string `toml:"commit_style"`
	AutoPush     *bool  `toml:"auto_push"`
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
	applyTOMLConfig(cfg, decoded, path)
	if err := validateConfig(cfg, path); err != nil {
		return err
	}
	return nil
}

func applyTOMLConfig(cfg *Config, decoded tomlConfig, path string) {
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
	if decoded.Fields.Lifecycle != "" {
		cfg.Fields.Lifecycle = decoded.Fields.Lifecycle
	}
	if decoded.Fields.RecordType != "" {
		cfg.Fields.RecordType = decoded.Fields.RecordType
	}
	if decoded.Fields.TaskValue != "" {
		cfg.Fields.TaskValue = decoded.Fields.TaskValue
	}
	if decoded.Fields.OutcomeValue != "" {
		cfg.Fields.OutcomeValue = decoded.Fields.OutcomeValue
	}
	if decoded.Fields.DisplayName != "" {
		cfg.Fields.DisplayName = decoded.Fields.DisplayName
	}
	if decoded.Fields.DependencyLink != "" {
		cfg.Fields.DependencyLink = decoded.Fields.DependencyLink
	}
	if decoded.Gitflow.BaseBranch != "" {
		cfg.Gitflow.BaseBranch = decoded.Gitflow.BaseBranch
	}
	if decoded.Gitflow.BranchMode != "" {
		cfg.Gitflow.BranchMode = decoded.Gitflow.BranchMode
	}
	if decoded.Gitflow.BranchStyle != "" {
		cfg.Gitflow.BranchStyle = decoded.Gitflow.BranchStyle
	}
	if decoded.Gitflow.PrCreate != "" {
		cfg.Gitflow.PrCreate = decoded.Gitflow.PrCreate
	}
	if decoded.Gitflow.PrTitleStyle != "" {
		cfg.Gitflow.PrTitleStyle = decoded.Gitflow.PrTitleStyle
	}
	if decoded.Gitflow.PrBodyStyle != "" {
		cfg.Gitflow.PrBodyStyle = decoded.Gitflow.PrBodyStyle
	}
	if decoded.Gitflow.CommitStyle != "" {
		cfg.Gitflow.CommitStyle = decoded.Gitflow.CommitStyle
	}
	if decoded.Gitflow.AutoPush != nil {
		cfg.Gitflow.AutoPush = *decoded.Gitflow.AutoPush
	}

	// Migration from deprecated top-level fields
	deprecatedPresent := decoded.CommitStyle != "" || decoded.AutoPush != nil || decoded.PRMode != nil
	if deprecatedPresent {
		cfg.Warnings = append(cfg.Warnings, Warning{
			Code:    "RMC_CONFIG_DEPRECATED_TOPLEVEL",
			Message: "top-level commit_style, auto_push, pr_mode are deprecated; move them under [gitflow]",
			Path:    path,
		})
		if decoded.CommitStyle != "" && decoded.Gitflow.CommitStyle == "" {
			cfg.Gitflow.CommitStyle = decoded.CommitStyle
		}
		if decoded.AutoPush != nil && decoded.Gitflow.AutoPush == nil {
			cfg.Gitflow.AutoPush = *decoded.AutoPush
		}
		if decoded.PRMode != nil && decoded.Gitflow.BranchMode == "" && decoded.Gitflow.PrCreate == "" {
			if *decoded.PRMode {
				cfg.Gitflow.BranchMode = "scope_branch"
				cfg.Gitflow.PrCreate = "auto"
			} else {
				cfg.Gitflow.BranchMode = "direct_push"
				cfg.Gitflow.PrCreate = "never"
			}
		}
	}
}


func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func roadmapRootExists(repo string, roadmapRoot string) bool {
	root, _, err := pathsec.ResolveInside(repo, roadmapRoot)
	if err != nil {
		return false
	}
	info, err := os.Stat(root)
	return err == nil && info.IsDir()
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
		Fields: FieldsConfig{
			Lifecycle:      "estado",
			RecordType:     "tipo",
			TaskValue:      "task",
			OutcomeValue:   "outcome",
			DisplayName:    "titulo",
			DependencyLink: "blocked_by",
		},
		LeafFilter:         "isIndex == false",
		OutcomeCloseVerify: []string{},
		Gitflow: GitflowConfig{
			BranchMode:  "direct_push",
			PrCreate:    "never",
			CommitStyle: "conventional",
			AutoPush:    true,
		},
		RequiredCodeCoverage:   85.0,
		LoopMaxTasks:           0,
		Parallel:               true,
		Autonomy:               "until_done",
		CompactAfterTaskCommit: true,
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
	default:
		return &Error{Code: ErrConfigParse, Message: "autonomy must be one of manual, supervised, until_done", Path: path, ExitCode: 2}
	}
	switch cfg.Gitflow.BranchMode {
	case "direct_push", "scope_branch":
	default:
		return &Error{Code: ErrConfigParse, Message: "gitflow.branch_mode must be one of direct_push, scope_branch", Path: path, ExitCode: 2}
	}
	switch cfg.Gitflow.PrCreate {
	case "never", "manual", "auto":
	default:
		return &Error{Code: ErrConfigParse, Message: "gitflow.pr_create must be one of never, manual, auto", Path: path, ExitCode: 2}
	}
	if cfg.Gitflow.PrCreate == "auto" && cfg.Gitflow.BranchMode != "scope_branch" {
		return &Error{Code: ErrConfigParse, Message: "gitflow.pr_create=auto requires gitflow.branch_mode=scope_branch", Path: path, ExitCode: 2}
	}
	return nil
}
