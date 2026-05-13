# Roadmap Skill Responsibility-Chain Redesign — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix three responsibility-chain violations: make `rootline describe` deterministic for multi-pattern schemas (adding `next_by_pattern`), make the LLM the sole writer of roadmap files (via Write tool), and remove the deprecated `materialize` and `plan-paths` CLI commands.

**Architecture:** Two repos are touched — `rootline` (Go code + skill docs) then `roadmapctl` (CLI deletion + skill rewrite). rootline changes must be committed first because the skill docs reference `next_by_pattern`. All Go changes follow TDD (failing test → implementation → green → commit).

**Tech Stack:** Go 1.22+, Cobra CLI, rootline schema system, Markdown skill files

---

## File Map

### rootline (`/home/shared/rootline`)

| File | Change |
|------|--------|
| `internal/rules/rules.go` | Add `NextByPattern` field to `SchemaField` struct |
| `internal/rules/describe.go` | Add `sortedPatterns`, fix `computeNextSequence`, add `computeAllNextSequences`, wire in `NewDescribeResult` |
| `internal/rules/compute_sequence_test.go` | Add tests for multi-pattern determinism and `computeAllNextSequences` |
| `internal/rules/describe_test.go` | Add test for `next_by_pattern` in JSON output |
| `.claude/skills/rootline/SKILL.md` | Add note about multi-pattern `next` + `next_by_pattern` |
| `.claude/skills/rootline/ref-schema.md` | Add `next_by_pattern` section |
| `docs/describe.md` | Add `next_by_pattern` example |

### roadmapctl (`/home/shared/roadmapctl`)

| File | Change |
|------|--------|
| `internal/cli/cli.go` | Remove lines 75–76 (`newMaterializeCommand`, `newPlanPathsCommand`) |
| `internal/cli/pathplan.go` | Delete |
| `internal/cli/pathplan_test.go` | Delete |
| `internal/cli/materialize.go` | Delete |
| `internal/materialize/` (6 files) | Delete package |
| `internal/roadmap/numbering.go` | Delete |
| `internal/roadmap/numbering_plan_test.go` | Delete |
| `docs/materialize-plan-schema.md` | Delete |
| `internal/cli/cli_test.go` | Remove `"materialize"` from `want` slice |
| `internal/cli/schema_compatibility_test.go` | Delete `TestMaterializeApplyBlocksStaleStemBeforeWritingPlannedFiles` |
| `docs/cli-contract.md` | Remove materialize command documentation |
| `.claude/skills/roadmap/SKILL.md` | Replace "Invariante de escritura segura" section |
| `.claude/skills/roadmap/plan-subcommand.md` | Full rewrite |
| `.claude/skills/roadmap/common-logic.md` | Surgical edits |
| `.claude/skills/roadmap/outcome-guide.md` | Update lines 16 and 42 |

---

## Task 1: Write failing tests for `next_by_pattern` (rootline)

**Repo:** `/home/shared/rootline`

**Files:**
- Modify: `internal/rules/compute_sequence_test.go`
- Modify: `internal/rules/describe_test.go`

- [ ] **Step 1: Add tests to compute_sequence_test.go**

Append to `internal/rules/compute_sequence_test.go` (after the existing `TestComputeNextSequence_ZeroDigits` test):

```go
func makeMultiPatternField(t *testing.T) (SchemaField, string) {
	t.Helper()
	dir := t.TempDir()
	mkTestDir(t, filepath.Join(dir, "O01-first-outcome"))
	mkTestDir(t, filepath.Join(dir, "O02-second-outcome"))
	writeTestFile(t, filepath.Join(dir, "T001-direct-task.md"))

	field := SchemaField{
		Type: "sequence",
		Match: &FieldMatch{
			Configs: map[string]any{
				"O*": map[string]any{"prefix": "O", "digits": 2},
				"T*": map[string]any{"prefix": "T", "digits": 3},
			},
		},
	}
	return field, dir
}

func TestComputeAllNextSequences_MultiPattern(t *testing.T) {
	field, dir := makeMultiPatternField(t)

	got := computeAllNextSequences(dir, field)

	if got == nil {
		t.Fatal("computeAllNextSequences returned nil, want map")
	}
	if got["O*"] != "O03" {
		t.Errorf("O*: got %q, want O03", got["O*"])
	}
	if got["T*"] != "T002" {
		t.Errorf("T*: got %q, want T002", got["T*"])
	}
}

func TestComputeNextSequence_MultiPatternDeterministic(t *testing.T) {
	field, dir := makeMultiPatternField(t)

	// O* < T* alphabetically; both patterns have matching entries.
	// next must always return the O* value — deterministically.
	first := computeNextSequence(dir, field)
	for i := 0; i < 10; i++ {
		got := computeNextSequence(dir, field)
		if got != first {
			t.Errorf("non-deterministic: run %d got %q, first run got %q", i+1, got, first)
		}
	}
	if first != "O03" {
		t.Errorf("next = %q, want O03 (first alphabetical pattern that matches)", first)
	}
}

func TestComputeAllNextSequences_NilMatch(t *testing.T) {
	dir := t.TempDir()
	field := SchemaField{Type: "sequence"}
	got := computeAllNextSequences(dir, field)
	if got != nil {
		t.Errorf("nil match: got %v, want nil", got)
	}
}
```

- [ ] **Step 2: Add test to describe_test.go**

Append to `internal/rules/describe_test.go` (after the last test):

```go
func TestDescribeResult_NextByPatternInJSON(t *testing.T) {
	dir := t.TempDir()
	// Create entries matching both O* and T* patterns
	if err := os.Mkdir(filepath.Join(dir, "O01-outcome"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(dir, "T001-task.md"))

	entries := []StemEntry{{Path: filepath.Join(dir, ".stem"), Stem: &StemFile{}}}
	effective := &StemFile{
		Schema: map[string]SchemaField{
			"id": {
				Type: "sequence",
				Match: &FieldMatch{
					Configs: map[string]any{
						"O*": map[string]any{"prefix": "O", "digits": 2},
						"T*": map[string]any{"prefix": "T", "digits": 3},
					},
				},
			},
		},
	}

	result := NewDescribeResult(dir, entries, effective)

	idField := result.Schema["id"]
	if idField.NextByPattern == nil {
		t.Fatal("NextByPattern is nil, want map")
	}
	if idField.NextByPattern["O*"] != "O02" {
		t.Errorf("NextByPattern[O*] = %q, want O02", idField.NextByPattern["O*"])
	}
	if idField.NextByPattern["T*"] != "T002" {
		t.Errorf("NextByPattern[T*] = %q, want T002", idField.NextByPattern["T*"])
	}

	// Verify it appears in JSON output
	data, err := result.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("JSON parse: %v", err)
	}
	schema := parsed["schema"].(map[string]any)
	idJSON := schema["id"].(map[string]any)
	nbp, ok := idJSON["next_by_pattern"].(map[string]any)
	if !ok {
		t.Fatalf("next_by_pattern missing from JSON or wrong type: %v", idJSON)
	}
	if nbp["O*"] != "O02" {
		t.Errorf("JSON next_by_pattern[O*] = %v, want O02", nbp["O*"])
	}
}
```

Note: `describe_test.go` already imports `"os"` and `"encoding/json"` — check top of file; add any missing imports.

- [ ] **Step 3: Run the new tests — verify they FAIL**

```bash
cd /home/shared/rootline
go test ./internal/rules/... -run "TestComputeAllNextSequences|TestComputeNextSequence_MultiPatternDeterministic|TestDescribeResult_NextByPatternInJSON" -v 2>&1 | head -30
```

Expected: compile error `undefined: computeAllNextSequences`. If tests pass without implementation, the test is wrong — fix it before continuing.

---

## Task 2: Add `NextByPattern` to `SchemaField` (rootline)

**Repo:** `/home/shared/rootline`

**Files:**
- Modify: `internal/rules/rules.go`

- [ ] **Step 1: Add the field after `Next string` (line 160)**

In `internal/rules/rules.go`, the `SchemaField` struct currently has:

```go
	Next          string       `yaml:"-" json:"next,omitempty"`
	Excludes      *ExcludeRule `yaml:"excludes" json:"excludes,omitempty"`
```

Change it to:

```go
	Next          string            `yaml:"-" json:"next,omitempty"`
	NextByPattern map[string]string `yaml:"-" json:"next_by_pattern,omitempty"`
	Excludes      *ExcludeRule      `yaml:"excludes" json:"excludes,omitempty"`
```

- [ ] **Step 2: Verify it compiles**

```bash
cd /home/shared/rootline
go build ./internal/rules/...
```

Expected: exit 0, no output.

---

## Task 3: Implement `computeAllNextSequences` + fix `computeNextSequence` (rootline)

**Repo:** `/home/shared/rootline`

**Files:**
- Modify: `internal/rules/describe.go`

- [ ] **Step 1: Add `"sort"` to imports**

In `internal/rules/describe.go`, the import block currently is:

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)
```

Add `"sort"`:

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)
```

- [ ] **Step 2: Replace `computeNextSequence` and add helpers**

Replace the entire `computeNextSequence` function (lines 104–140) with:

```go
// computeNextSequence scans dirPath for files/dirs matching the prefix pattern,
// finds the highest numeric suffix, and returns prefix + next number zero-padded
// to the specified digits. Supports both top-level prefix/digits and match configs
// where prefix/digits are nested per-pattern (e.g., "E*": {prefix: E, digits: 2}).
// When multiple match config patterns are present, patterns are evaluated in
// alphabetical order to ensure deterministic results.
func computeNextSequence(dirPath string, field SchemaField) string {
	if field.Prefix != "" && field.Digits > 0 {
		return computeNextFromPrefix(dirPath, field.Prefix, field.Digits)
	}

	if field.Match == nil || field.Match.Configs == nil {
		return ""
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return ""
	}

	for _, globPattern := range sortedPatterns(field.Match.Configs) {
		config := field.Match.Configs[globPattern]
		cfgMap, ok := config.(map[string]any)
		if !ok {
			continue
		}
		prefix, digits := extractPrefixDigits(cfgMap)
		if prefix == "" || digits <= 0 {
			continue
		}

		for _, e := range entries {
			if matched, _ := filepath.Match(globPattern, e.Name()); matched {
				return computeNextFromPrefix(dirPath, prefix, digits)
			}
		}
	}

	return ""
}

// computeAllNextSequences returns the next sequence value for every pattern
// in match configs. Unlike computeNextSequence, it does not stop at the first
// match — it computes a value for each pattern independently.
func computeAllNextSequences(dirPath string, field SchemaField) map[string]string {
	if field.Match == nil || field.Match.Configs == nil {
		return nil
	}

	result := make(map[string]string, len(field.Match.Configs))
	for _, globPattern := range sortedPatterns(field.Match.Configs) {
		config := field.Match.Configs[globPattern]
		cfgMap, ok := config.(map[string]any)
		if !ok {
			continue
		}
		prefix, digits := extractPrefixDigits(cfgMap)
		if prefix == "" || digits <= 0 {
			continue
		}
		result[globPattern] = computeNextFromPrefix(dirPath, prefix, digits)
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// sortedPatterns returns the keys of configs sorted alphabetically.
func sortedPatterns(configs map[string]any) []string {
	patterns := make([]string, 0, len(configs))
	for p := range configs {
		patterns = append(patterns, p)
	}
	sort.Strings(patterns)
	return patterns
}
```

- [ ] **Step 3: Wire `NextByPattern` into `NewDescribeResult`**

In `NewDescribeResult`, find the sequence field block (lines 59–65):

```go
	// Compute Next for sequence fields
	for name, field := range schema {
		if field.Type == "sequence" {
			field.Next = computeNextSequence(path, field)
			schema[name] = field
		}
	}
```

Replace it with:

```go
	// Compute Next and NextByPattern for sequence fields
	for name, field := range schema {
		if field.Type == "sequence" {
			field.Next = computeNextSequence(path, field)
			field.NextByPattern = computeAllNextSequences(path, field)
			schema[name] = field
		}
	}
```

- [ ] **Step 4: Run all tests — verify they PASS**

```bash
cd /home/shared/rootline
go test ./internal/rules/... -v 2>&1 | tail -20
```

Expected: all PASS. If any existing test regresses, fix describe.go before continuing.

- [ ] **Step 5: Commit**

```bash
cd /home/shared/rootline
git add internal/rules/rules.go internal/rules/describe.go internal/rules/compute_sequence_test.go internal/rules/describe_test.go
git commit -m "feat(describe): add next_by_pattern for deterministic multi-pattern sequence numbering"
```

---

## Task 4: Remove `plan-paths` CLI from roadmapctl

**Repo:** `/home/shared/roadmapctl`

**Files:**
- Delete: `internal/cli/pathplan.go`
- Delete: `internal/cli/pathplan_test.go`
- Modify: `internal/cli/cli.go`

- [ ] **Step 1: Delete pathplan files**

```bash
cd /home/shared/roadmapctl
rm internal/cli/pathplan.go internal/cli/pathplan_test.go
```

- [ ] **Step 2: Remove `newPlanPathsCommand` from `cli.go`**

In `internal/cli/cli.go`, remove line 76:

```go
	cmd.AddCommand(newPlanPathsCommand(&options, stdout, stderr, exitCode))
```

The surrounding context (keep these):

```go
	cmd.AddCommand(newMaterializeCommand(&options, stdout, stderr, exitCode))
	// DELETE the line above ↑ in Task 5; only delete plan-paths here
	cmd.AddCommand(newBootstrapCommand(&options, stdout, stderr, exitCode))
```

After edit, line 76 (previously plan-paths) no longer exists; `newBootstrapCommand` follows `newMaterializeCommand`.

- [ ] **Step 3: Build to confirm no compile errors**

```bash
cd /home/shared/roadmapctl
go build ./...
```

Expected: exit 0.

- [ ] **Step 4: Commit**

```bash
cd /home/shared/roadmapctl
git add internal/cli/cli.go
git rm internal/cli/pathplan.go internal/cli/pathplan_test.go
git commit -m "feat(cli): remove plan-paths command (replaced by rootline next_by_pattern)"
```

---

## Task 5: Remove `materialize` CLI and package from roadmapctl

**Repo:** `/home/shared/roadmapctl`

**Files:**
- Modify: `internal/cli/cli.go`
- Delete: `internal/cli/materialize.go`
- Delete: `internal/materialize/` (6 files)
- Delete: `internal/roadmap/numbering.go`
- Delete: `internal/roadmap/numbering_plan_test.go`
- Modify: `internal/cli/cli_test.go`
- Modify: `internal/cli/schema_compatibility_test.go`

- [ ] **Step 1: Remove `newMaterializeCommand` from `cli.go`**

In `internal/cli/cli.go`, remove the line:

```go
	cmd.AddCommand(newMaterializeCommand(&options, stdout, stderr, exitCode))
```

After this edit, `newBootstrapCommand` immediately follows `newTransitionCommand`.

- [ ] **Step 2: Delete materialize CLI and package files**

```bash
cd /home/shared/roadmapctl
rm internal/cli/materialize.go
rm internal/materialize/dryrun.go \
   internal/materialize/pathplan.go \
   internal/materialize/dryrun_test.go \
   internal/materialize/pathplan_test.go \
   internal/materialize/apply_changes_test.go \
   internal/materialize/validation_test.go
rmdir internal/materialize
rm internal/roadmap/numbering.go internal/roadmap/numbering_plan_test.go
```

- [ ] **Step 3: Update `cli_test.go` — remove "materialize" from help check**

In `internal/cli/cli_test.go`, `TestHelpListsImplementedCommands` has:

```go
	for _, want := range []string{"roadmapctl", "doctor", "check", "materialize", "--output", "--repo"} {
```

Change to:

```go
	for _, want := range []string{"roadmapctl", "doctor", "check", "--output", "--repo"} {
```

- [ ] **Step 4: Delete materialize test from `schema_compatibility_test.go`**

In `internal/cli/schema_compatibility_test.go`, delete the entire function `TestMaterializeApplyBlocksStaleStemBeforeWritingPlannedFiles` (lines 58–73):

```go
func TestMaterializeApplyBlocksStaleStemBeforeWritingPlannedFiles(t *testing.T) {
	fixture := copyFixture(t, "valid-outcome-with-tasks")
	writeStaleOutcomeEstadoStem(t, fixture)
	before := listRoadmapFiles(t, fixture)
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"materialize", "--plan", plan, "--apply", "--repo", fixture, "--output", "json"}, &stdout, &stderr, "dev")
	testutil.AssertExit(t, code, 1, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	testutil.RequireDiagnosticID(t, report, "RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED")
	after := listRoadmapFiles(t, fixture)
	if before != after {
		t.Fatalf("materialize wrote files despite stale stem\nbefore:\n%s\nafter:\n%s", before, after)
	}
}
```

After deletion, check if `filepath` is still imported elsewhere in that file. If the deleted function was the only use, remove the `"path/filepath"` import too.

- [ ] **Step 5: Build and run tests**

```bash
cd /home/shared/roadmapctl
go build ./...
go test ./... 2>&1 | tail -20
```

Expected: build succeeds, all tests pass. If any test references a deleted function, fix the test.

- [ ] **Step 6: Commit**

```bash
cd /home/shared/roadmapctl
git add internal/cli/cli.go internal/cli/cli_test.go internal/cli/schema_compatibility_test.go
git rm internal/cli/materialize.go \
       internal/materialize/dryrun.go \
       internal/materialize/pathplan.go \
       internal/materialize/dryrun_test.go \
       internal/materialize/pathplan_test.go \
       internal/materialize/apply_changes_test.go \
       internal/materialize/validation_test.go \
       internal/roadmap/numbering.go \
       internal/roadmap/numbering_plan_test.go
git commit -m "feat(cli): remove materialize command and package (deprecated by ADR 2026-05-11)"
```

---

## Task 6: Update `docs/cli-contract.md`

**Repo:** `/home/shared/roadmapctl`

**Files:**
- Modify: `docs/cli-contract.md`
- Delete: `docs/materialize-plan-schema.md`

- [ ] **Step 1: Remove materialize from the command list (line 17)**

Find and remove the line:
```
- `roadmapctl materialize`
```

- [ ] **Step 2: Remove materialize from the help output example (around line 47)**

Find and remove the line:
```
  materialize  Validate and write approved structured roadmap plans.
```

- [ ] **Step 3: Remove "materialize" from read-only guarantees (lines 184, 207)**

Line 184 — change:
```
`check` must not write, materialize, fix, or normalize roadmap files.
```
To:
```
`check` must not write, fix, or normalize roadmap files.
```

Line 207 — change:
```
it must not materialize, normalize, auto-fix, or judge subjective writing quality.
```
To:
```
it must not normalize, auto-fix, or judge subjective writing quality.
```

- [ ] **Step 4: Remove materialize from diagnostic table (lines 392–393)**

In the diagnostic table, remove `, materialize` from the "Commands" column of `RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED` and `RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY` rows.

- [ ] **Step 5: Remove materialize command documentation block (around lines 456–466)**

Find and remove the paragraph starting with:
```
`roadmapctl materialize` accepts only structured, approved JSON input.
```
...through the end of the code block showing `roadmapctl materialize ...` invocations.

- [ ] **Step 6: Delete materialize schema doc**

```bash
cd /home/shared/roadmapctl
rm docs/materialize-plan-schema.md
```

- [ ] **Step 7: Commit**

```bash
cd /home/shared/roadmapctl
git add docs/cli-contract.md
git rm docs/materialize-plan-schema.md
git commit -m "docs(cli-contract): remove materialize command references"
```

---

## Task 7: Update `SKILL.md` write invariant

**Repo:** `/home/shared/roadmapctl`

**Files:**
- Modify: `.claude/skills/roadmap/SKILL.md`

- [ ] **Step 1: Replace the "Invariante de escritura segura" section**

Find the section starting at line ~70:

```
## Invariante de escritura segura

Además del criterio canónico, el skill no puede crear/reescribir archivos roadmap manualmente ni hacer dumps multi-file fuera de `roadmapctl`.

Prohibido:

- `bash` con múltiples heredocs/cats dirigidos a archivos distintos.
- loops de shell que llamen `rootline new` o escritura en varios paths.
- `write`/`edit` directo para crear o reescribir archivos canónicos del roadmap.

Permitido: `roadmapctl materialize --plan <plan-json> --apply` o `roadmapctl materialize --changes <dry-run-json> --apply` puede aplicar múltiples archivos en una ejecución, porque roadmapctl owns canonical writes, per-file diagnostics, validation, ordering and postcheck.
```

Replace it with:

```
## Invariante de escritura segura

El skill es el único writer de archivos roadmap, vía Write tool, después de aprobación explícita del usuario y preflight pasado.

Prohibido:

- `bash`/`sh` con múltiples heredocs, `cat >`, o loops que escriban varios archivos en una sola llamada.
- loops de shell que llamen `rootline new` para múltiples paths.
- Escribir archivos roadmap sin que el usuario haya aprobado el árbol propuesto.
- Escribir archivos roadmap si `roadmapctl doctor` o `roadmapctl check --strict` retornan non-zero.

Permitido: Write tool para cada archivo canónico (`OXX-slug/README.md`, `OXX-slug/TXXX-slug.md`) después de aprobación y preflight exitoso.
```

- [ ] **Step 2: Commit**

```bash
cd /home/shared/roadmapctl
git add .claude/skills/roadmap/SKILL.md
git commit -m "docs(skill): update write invariant — LLM writes directly via Write tool"
```

---

## Task 8: Rewrite `plan-subcommand.md`

**Repo:** `/home/shared/roadmapctl`

**Files:**
- Modify: `.claude/skills/roadmap/plan-subcommand.md`

- [ ] **Step 1: Replace the entire file content**

Write `.claude/skills/roadmap/plan-subcommand.md` with:

```markdown
# /roadmap plan

Materializa el plan de la conversación como archivos `.md` del roadmap. No implementa código.

Ruta normal autosuficiente: este archivo contiene el procedimiento operativo completo. No leer `common-logic.md` ni documentación de integración para ejecutar el flujo; esos documentos son referencia de mantenimiento/troubleshooting.

Materializar es una operación estructural. Está prohibido crear un único archivo
con una lista de tareas. Cada task debe tener su propio archivo `TXXX-*.md`.

## Fuente del plan

1. Contexto actual de conversación.

Si no hay plan, informar: "No hay plan en esta conversación. Primero planificar, luego ejecutar `/roadmap plan`." y parar.

## Workspace mode

Resolver repo target:

1. `--repo <name>` si fue dado.
2. Repo mencionado en el plan.
3. Si ambiguo, preguntar.

Usar `<abs-roadmap-root>` y `git -C <repo-path>`.

## Fase 1: Descomposición

1. Identificar el plan más reciente de la conversación.

2. Consultar numeración actual:
   ```bash
   rootline describe <roadmap-root> --field schema.id.next_by_pattern --output json
   ```
   Retorna `{"O*": "O14", "T*": "T014"}`.
   - Usar `O*` para el siguiente Outcome.
   - Usar `T*` como referencia inicial para tasks en outcomes nuevos.

   Para tasks en un **Outcome existente**:
   ```bash
   rootline describe <roadmap-root>/OXX-slug/ --field schema.id.next_by_pattern --output json
   ```
   Retorna `{"T*": "T009"}` — primer task disponible dentro de ese outcome.

   Para tasks en un **Outcome nuevo** (directorio aún no existe): comenzar tasks desde T001.

3. Aplicar `framework-reference.md`: máximo Outcome + Tasks por outcome.
4. Asignar slugs (kebab-case, sin prefijos O/T, sin `/` ni `..`) y numerar con los valores obtenidos.
5. Cada task: nombre, descripción, ACs principales, `hard_blockers` solo si hay dependencia objetiva real.

## Fase 2: Aprobación

Presentar árbol completo con números asignados + ACs:

```
O14-nombre-outcome/
├── README.md
├── T001-primera-task.md
│   - AC1: ...
└── T002-segunda-task.md
    - AC1: ...
```

**STOP obligatorio** con `AskUserQuestion` hasta aprobación explícita. No crear archivos antes.

## Fase 3: Materialización

**3.1 Re-confirmar numeración (antistaleness)**

```bash
rootline describe <roadmap-root> --field schema.id.next_by_pattern --output json
```

Si aparecieron nuevos archivos que cambian los números propuestos, informar al usuario y recalcular antes de continuar.

**3.2 Preflight obligatorio**

```bash
command -v roadmapctl
roadmapctl doctor --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si cualquier comando sale non-zero: detenerse, reportar exit code y diagnostics. No crear archivos.

**3.3 Escritura en paralelo**

Crear directorios padre si aplican, luego escribir con Write tool en paralelo:

- `OXX-slug/README.md`: frontmatter `tipo: outcome` + título + descripción/contexto (SIN `## Criterios de Aceptación` ni `## Tasks`). Ver template en `outcome-guide.md`.
- `OXX-slug/TXXX-slug.md`: frontmatter `estado: Specified` + título + descripción + `## Criterios de Aceptación` + contexto + scope + hard blockers si aplican. Ver template en `task-guide.md`.

**3.4 Validación por archivo**

```bash
rootline validate <path-creado>
```

Por cada archivo creado. Si falla: reportar y detener.

**3.5 Postcheck obligatorio**

```bash
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si falla: detenerse, reportar diagnostics. No commitear.

## Fase 4: Commit

```bash
git -C <repo-path> add <archivos .md creados>
git -C <repo-path> commit -m "chore(roadmap): create planning docs"
```

STOP. Informar: "Archivos de planificación creados. Ejecutar `/roadmap loop` cuando esté listo para implementar."
```

- [ ] **Step 2: Commit**

```bash
cd /home/shared/roadmapctl
git add .claude/skills/roadmap/plan-subcommand.md
git commit -m "docs(skill): rewrite plan-subcommand — LLM writes directly, rootline next_by_pattern for numbering"
```

---

## Task 9: Surgical edits to `common-logic.md`

**Repo:** `/home/shared/roadmapctl`

**Files:**
- Modify: `.claude/skills/roadmap/common-logic.md`

- [ ] **Step 1: Delete "Invariante de escritura segura" block**

Find and delete the entire section from `## Invariante de escritura segura` through the closing paragraph that ends with `...corre postcheck antes de éxito.` (the paragraph ending with `porque `roadmapctl` valida el plan/change-set, ordena padres antes de hijos, reporta diagnostics por path, ejecuta validaciones y corre postcheck antes de éxito.`). This includes the `**Prohibido en una misma tool call del skill:**` list and the `Permitido:` paragraph.

- [ ] **Step 2: Delete "Materialización determinística" section**

Find and delete the section starting with `## Materialización determinística` through the end of the bash block showing `roadmapctl materialize --plan ... --apply`.

- [ ] **Step 3: Delete "Verificación de padre" section**

Find and delete the section starting with `## Verificación de padre` through the sentence ending with `No ejecutar `rootline describe` como paso primario antes de crear archivos.`

- [ ] **Step 4: Replace "Auto-numbering" section**

Find:
```
## Auto-numbering

El skill no calcula números `OXX`/`TXXX`. `roadmapctl materialize` asigna numbering determinístico y reporta las rutas propuestas en `changes[]` durante dry-run. Si el dry-run no produce rutas canónicas, detenerse y reportar diagnostics.
```

Replace with:
```
## Auto-numbering

El skill usa `rootline describe` para obtener numeración determinística:

```bash
# Retorna el siguiente O y T simultáneamente
rootline describe <roadmap-root> --field schema.id.next_by_pattern --output json
# → {"O*": "O14", "T*": "T014"}

# Retorna solo el siguiente T dentro de un outcome existente
rootline describe <roadmap-root>/OXX-slug/ --field schema.id.next_by_pattern --output json
# → {"T*": "T009"}
```

Preferir `next_by_pattern` (mapa) sobre `next` (string) en schemas con múltiples patrones de secuencia. `next` es determinístico post-fix pero retorna solo el primer patrón alfabético que coincide con entries existentes.
```

- [ ] **Step 5: Commit**

```bash
cd /home/shared/roadmapctl
git add .claude/skills/roadmap/common-logic.md
git commit -m "docs(skill): remove materialize/invariant sections, update auto-numbering to rootline next_by_pattern"
```

---

## Task 10: Update `outcome-guide.md`

**Repo:** `/home/shared/roadmapctl`

**Files:**
- Modify: `.claude/skills/roadmap/outcome-guide.md`

- [ ] **Step 1: Update line 16**

Find:
```
El skill no calcula `OXX`. `roadmapctl path-planning` (cuando esté disponible) asigna el siguiente Outcome determinísticamente y muestra la ruta propuesta. Mientras tanto, usar `roadmapctl next` como referencia para numbering determinístico.
```

Replace with:
```
El skill no calcula `OXX` manualmente. Usar `rootline describe <roadmap-root> --field schema.id.next_by_pattern --output json` para obtener el siguiente O y T determinísticamente.
```

- [ ] **Step 2: Update line 42**

Find:
```
- El title del README usa el título del plan (`# [Nombre del objetivo]`); el identificador `OXX` vive en la ruta asignada por path-planning.
```

Replace with:
```
- El title del README usa el título del plan (`# [Nombre del objetivo]`); el identificador `OXX` viene del valor `O*` retornado por `rootline describe --field schema.id.next_by_pattern`.
```

- [ ] **Step 3: Commit**

```bash
cd /home/shared/roadmapctl
git add .claude/skills/roadmap/outcome-guide.md
git commit -m "docs(skill): update outcome-guide — rootline next_by_pattern replaces path-planning"
```

---

## Task 11: Update rootline skill docs

**Repo:** `/home/shared/rootline`

**Files:**
- Modify: `.claude/skills/rootline/SKILL.md`
- Modify: `.claude/skills/rootline/ref-schema.md`
- Modify: `docs/describe.md`

- [ ] **Step 1: Update SKILL.md — add multi-pattern note**

In `.claude/skills/rootline/SKILL.md`, find the section that documents `schema.id.next` (around line 87). After the existing description, add:

```
> **Multi-pattern schemas:** `schema.id.next` retorna el próximo valor del primer patrón alfabético que tiene entries existentes en el directorio. En schemas con múltiples patrones de secuencia (ej: `O*` y `T*`), usar `--field schema.id.next_by_pattern` para obtener el próximo valor de **todos** los patrones simultáneamente: `{"O*": "O14", "T*": "T014"}`.
```

- [ ] **Step 2: Update ref-schema.md — add `next_by_pattern` section**

In `.claude/skills/rootline/ref-schema.md`, find where `schema.id.next` is documented (around lines 12, 84, 93). After the `next` field description, add:

```markdown
### Schemas multi-patrón: `next` vs `next_by_pattern`

Cuando `id` define múltiples patrones vía `match` (ej: `O*` y `T*`), el campo `next` retorna el próximo valor del primer patrón alfabético que coincide con entries existentes en el directorio — determinístico, pero incompleto para schemas multi-patrón.

Para obtener el próximo valor de **cada** patrón:

```bash
rootline describe <dir> --field schema.id.next_by_pattern
# → {"O*": "O14", "T*": "T014"}

rootline describe <dir>/O14-slug/ --field schema.id.next_by_pattern
# → {"T*": "T001"}
```

Usar `next_by_pattern` cuando el LLM necesita asignar números tanto para Outcomes como para Tasks en el mismo flujo de materialización.
```

- [ ] **Step 3: Update docs/describe.md — add example**

In `docs/describe.md`, in the section documenting `--field`, add after the existing `schema.id.next` example:

```markdown
**Multi-pattern sequences** — get the next value for each pattern:

```bash
rootline describe docs/roadmap/ --field schema.id.next_by_pattern
# → {"O*": "O14", "T*": "T014"}
```
```

- [ ] **Step 4: Commit**

```bash
cd /home/shared/rootline
git add .claude/skills/rootline/SKILL.md .claude/skills/rootline/ref-schema.md docs/describe.md
git commit -m "docs(skill): document next_by_pattern for multi-pattern sequence schemas"
```

---

## Verification Checklist

After all tasks are complete:

- [ ] `go test ./internal/rules/... -v` in rootline — all PASS
- [ ] `rootline describe <roadmap-root> --field schema.id.next_by_pattern` returns `{"O*":"OXX","T*":"TXXX"}` — run 5 times consecutively and verify same output every time
- [ ] `roadmapctl --help` — does NOT show `materialize` or `plan-paths`
- [ ] `go build ./...` in roadmapctl — exit 0
- [ ] `go test ./...` in roadmapctl — all PASS
- [ ] `/roadmap plan` flow executes without calling `roadmapctl plan-paths` or `roadmapctl materialize`
