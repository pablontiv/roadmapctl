package templates

import (
	"strings"
	"testing"
)

func TestDefaultRoadmapctlTOMLIncludesRequiredCodeCoverage(t *testing.T) {
	if want := "required_code_coverage = 85.0"; !strings.Contains(DefaultRoadmapctlTOML, want) {
		t.Fatalf("DefaultRoadmapctlTOML missing %q:\n%s", want, DefaultRoadmapctlTOML)
	}
}

func TestGenerateStemContent(t *testing.T) {
	tests := []struct {
		name           string
		dependencyLink string
		wantContains   []string
	}{
		{
			name:           "blocked_by link",
			dependencyLink: "blocked_by",
			wantContains: []string{
				"version: 2",
				"scope:",
				"schema:",
				"estado:",
				"tipo:",
				"id:",
				"links:",
				"blocked_by:",
				"reference:",
				"validate:",
				"T[0-9]{3}",
				"values: [Pending, Specified, In Progress, Completed, Blocked, On Hold, Obsolete]",
				"values: [outcome, task]",
				"prefix: O, digits: 2",
				"prefix: T, digits: 3",
				"rule: non_empty",
			},
		},
		{
			name:           "custom_depends_on link",
			dependencyLink: "custom_depends_on",
			wantContains: []string{
				"custom_depends_on:",
				"T[0-9]{3}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := GenerateStemContent(tt.dependencyLink)

			// Verify it's non-empty
			if content == "" {
				t.Fatal("GenerateStemContent returned empty string")
			}

			// Verify all expected strings are present
			for _, want := range tt.wantContains {
				if !strings.Contains(content, want) {
					t.Errorf("GenerateStemContent(%q) missing %q\nfull content:\n%s", tt.dependencyLink, want, content)
				}
			}

			// Verify the dependency link is interpolated correctly
			if !strings.Contains(content, tt.dependencyLink+":") {
				t.Errorf("GenerateStemContent(%q) didn't interpolate dependency link", tt.dependencyLink)
			}
		})
	}
}

func TestGenerateStemContentWithDifferentLinks(t *testing.T) {
	links := []string{"blocked_by", "depends_on", "custom_field", "ref"}

	for _, link := range links {
		t.Run("link="+link, func(t *testing.T) {
			content := GenerateStemContent(link)

			// Verify version and scope are present
			if !strings.Contains(content, "version: 2") {
				t.Error("missing version: 2")
			}
			if !strings.Contains(content, `scope:`) {
				t.Error("missing scope")
			}

			// Verify the link is in the links section
			if !strings.Contains(content, link+":") {
				t.Errorf("link %q not found in content", link)
			}

			// Verify required schema fields are present
			requiredFields := []string{"estado:", "tipo:", "id:"}
			for _, field := range requiredFields {
				if !strings.Contains(content, field) {
					t.Errorf("required field %q missing", field)
				}
			}
		})
	}
}

func TestDefaultRoadmapctlTOMLIncludesAllRequiredSettings(t *testing.T) {
	requiredSettings := []string{
		"done_statuses",
		"active_statuses",
		"leaf_filter",
		"outcome_close_verify",
		"required_code_coverage = 85.0",
		"loop_max_tasks",
		"parallel",
		"autonomy",
		"compact_after_task_commit",
		"[gitflow]",
		`branch_mode = "direct_push"`,
		`pr_create = "never"`,
		`commit_style = "conventional"`,
		"auto_push = true",
		"[status_values]",
		"pending",
		"specified",
		"in_progress",
		"completed",
		"blocked",
		"obsolete",
	}

	for _, setting := range requiredSettings {
		if !strings.Contains(DefaultRoadmapctlTOML, setting) {
			t.Errorf("DefaultRoadmapctlTOML missing %q", setting)
		}
	}

	// Deprecated fields should not appear in the template
	for _, deprecated := range []string{"pr_mode", "commit_style = "} {
		// commit_style should only appear inside [gitflow], not before it at top level
		gitflowIdx := strings.Index(DefaultRoadmapctlTOML, "[gitflow]")
		if strings.Contains(DefaultRoadmapctlTOML[:gitflowIdx], deprecated) {
			t.Errorf("DefaultRoadmapctlTOML has deprecated %q before [gitflow] section", deprecated)
		}
	}
	if strings.Contains(DefaultRoadmapctlTOML, "pr_mode") {
		t.Error("DefaultRoadmapctlTOML should not contain deprecated pr_mode")
	}
}
