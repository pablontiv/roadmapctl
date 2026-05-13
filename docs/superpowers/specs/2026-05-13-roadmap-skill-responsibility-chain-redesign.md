# Design: Roadmap Skill Responsibility-Chain Redesign

Date: 2026-05-13  
Status: Approved  
ADR reference: `docs/decisions/materialize-writer-vs-guard-flow.md` (2026-05-11)

---

## Problem

The `/roadmap plan` skill violates the established responsibility chain:

```
llm → skill → roadmapctl → rootline → *.md
```

Three concrete violations:

1. **`roadmapctl materialize --apply` as writer** — The skill delegates file creation to `materialize`, which regenerates content the LLM already produced. This contradicts the ADR that designates the skill (Pi/LLM) as sole writer after human approval.

2. **"Invariante de escritura segura" forbids Write tool** — `common-logic.md` and `SKILL.md` prohibit the LLM from using Write tool for roadmap files. This prohibition directly contradicts the ADR.

3. **`rootline describe --field schema.id.next` is non-deterministic** — For schemas with multiple sequence patterns (e.g., `O*` and `T*`), Go map iteration randomizes which pattern is evaluated first, returning `"O14"` or `"T014"` unpredictably. This made `roadmapctl plan-paths` a necessary workaround; with a reliable `next_by_pattern` field, plan-paths becomes redundant.

---

## Design

### Responsibility boundary (final)

| Layer | Responsibility |
|-------|---------------|
| LLM | Decompose plan, assign slugs, generate file content, write files with Write tool |
| skill | Orchestrate: consult rootline for numbering, prompt user for approval, run preflight/postcheck |
| roadmapctl | Guard/policy: `doctor`, `check`, `bootstrap`; no file generation |
| rootline | Schema, validation, sequence numbering; no roadmap-specific knowledge |

### Change 1: Fix rootline — deterministic `next_by_pattern`

**File:** `internal/rules/describe.go`

`computeNextSequence` iterates `field.Match.Configs` (a Go map) non-deterministically. Fix: sort pattern keys alphabetically before iterating. Add `computeAllNextSequences` that returns the next sequence value for **every** pattern, not just the first match.

**File:** `internal/rules/rules.go`

Add `NextByPattern map[string]string` field to `SchemaField` (JSON: `next_by_pattern`, yaml: `-`). Populate it in `NewDescribeResult` alongside the existing `Next` field.

**Result:**

```bash
rootline describe <roadmap-root> --field schema.id.next_by_pattern
# → {"O*": "O14", "T*": "T014"}  — deterministic, always

rootline describe <roadmap-root>/O14-slug/ --field schema.id.next_by_pattern
# → {"T*": "T001"}

rootline describe <roadmap-root> --field schema.id.next
# → "O14"  — deterministic post-fix (first alphabetical pattern)
```

**Tests** (in `internal/rules/compute_sequence_test.go` + `describe_test.go`):
- Multi-pattern schema with both O* and T* entries → `next_by_pattern` returns both keys
- `computeNextSequence` is deterministic across 10+ runs on the same directory
- `nil` Match → `computeAllNextSequences` returns nil

### Change 2: Rewrite `plan-subcommand.md`

The skill no longer calls `roadmapctl plan-paths` or `roadmapctl materialize`. New flow:

**Phase 1 — Decomposition**
1. Read plan from conversation context.
2. Query numbering: `rootline describe <roadmap-root> --field schema.id.next_by_pattern --output json` → `{"O*": "O14", "T*": "T014"}`.  
   For tasks in an existing outcome: same command on the outcome dir → `{"T*": "T009"}`.  
   For tasks in a new outcome: start tasks from T001.
3. LLM assigns O/T numbers from the queried values, generates kebab-case slugs.
4. Each task gets: name, description, ACs, hard blockers only where objectively required.

**Phase 2 — Approval**  
Present numbered tree with ACs. STOP with `AskUserQuestion`. No files created before approval.

**Phase 3 — Materialization**
1. Re-query `next_by_pattern` (antistaleness check). If numbers changed, inform user.
2. Preflight: `roadmapctl doctor` + `roadmapctl check --strict`.
3. Write files in parallel with Write tool:
   - `OXX-slug/README.md` — frontmatter `tipo: outcome`, narrative context, no AC/Tasks sections.
   - `OXX-slug/TXXX-slug.md` — frontmatter `estado: Specified`, ACs in `## Criterios de Aceptación`.
4. `rootline validate <path>` per file. Stop on failure.
5. Postcheck: `roadmapctl check --strict`. No commit until this passes.

**Phase 4 — Commit**  
`git add` only the created `.md` files, then commit.

### Change 3: Surgical edits to `common-logic.md`

**Remove entirely:**
- "Invariante de escritura segura" block — prohibited Write tool (ADR violation)
- "Materialización determinística" section — primary path via `roadmapctl materialize`
- "Verificación de padre" — referenced `materialize --dry-run`
- "Auto-numbering" section (replaced by updated version below)

**Add/update:**
- Auto-numbering: "Use `rootline describe <path> --field schema.id.next_by_pattern`. Root returns next O and T simultaneously. Outcome dir returns next T only. Prefer `next_by_pattern` over `next` for multi-pattern schemas."

**Keep unchanged:** doctor/check guard, `rootline validate` post-write, no-fallback prohibition, rootline reference table, cascading links, `blocked_by` rules.

### Change 4: Remove `roadmapctl plan-paths` CLI

Delete from surface:
- `internal/cli/pathplan.go`
- `internal/cli/pathplan_test.go`
- `internal/cli/cli.go:76` — remove `newPlanPathsCommand`

Internals (`internal/materialize/pathplan.go`, `internal/roadmap/numbering.go`) removed together with materialize in Change 7.

### Change 5: Update `SKILL.md` invariant

Replace "Invariante de escritura segura" (lines 70–80):
- **Remove:** prohibition on Write tool for canonical roadmap files
- **Keep:** prohibition on bash multi-file heredoc dumps, uncontrolled `rootline new` loops
- **Keep:** requirement for user approval + preflight before any write
- **Remove:** all mentions of `roadmapctl materialize` (command no longer exists)

### Change 6: Update `outcome-guide.md`

Lines 16 and 42: replace "roadmapctl path-planning (cuando esté disponible)" with `rootline describe --field schema.id.next_by_pattern`.

### Change 7: Remove `roadmapctl materialize` CLI

The ADR (2026-05-11) deprecated `materialize`; the command was never actually deleted.

**Delete entirely:**
- `internal/cli/materialize.go`
- `internal/materialize/` package (6 files: `dryrun.go`, `pathplan.go`, `dryrun_test.go`, `pathplan_test.go`, `apply_changes_test.go`, `validation_test.go`)
- `internal/roadmap/numbering.go` + `numbering_plan_test.go`
- `docs/materialize-plan-schema.md`

**Remove from `cli.go`:** line 75 (`newMaterializeCommand`)

**Update:**
- `docs/cli-contract.md` — remove materialize section (lines 17, 47, 184, 207, 392+)
- `internal/cli/cli_test.go` — remove materialize subcommand tests
- `internal/cli/schema_compatibility_test.go` — remove materialize references
- `internal/roadmap/structure.go` — investigate and remove materialize reference

### Change 8: Update rootline skill docs

**`rootline/.claude/skills/rootline/SKILL.md`** (line ~87): Add note that `schema.id.next` is deterministic post-fix but incomplete for multi-pattern schemas; use `next_by_pattern` to get all patterns.

**`rootline/.claude/skills/rootline/ref-schema.md`**: Add section "Schemas multi-patrón: next vs next_by_pattern" with usage example.

**`rootline/docs/describe.md`**: Add `next_by_pattern` example in `--field` section.

These docs changes land after the rootline code is merged.

---

## Implementation Order

1. rootline code (tests first, then implementation) → build + test
2. roadmapctl: delete plan-paths CLI → delete materialize CLI + package → update cli.go, docs
3. roadmapctl skills: rewrite plan-subcommand.md, common-logic.md, SKILL.md, outcome-guide.md
4. rootline skills: SKILL.md, ref-schema.md, describe.md

---

## Verification

1. `rootline describe <roadmap-root> --field schema.id.next_by_pattern` returns `{"O*":"O14","T*":"T014"}` deterministically in 10+ consecutive runs.
2. `roadmapctl --help` shows no `materialize` or `plan-paths` subcommands.
3. `go build ./...` passes in both rootline and roadmapctl after deletions.
4. `/roadmap plan` flow: LLM calls `rootline describe` for numbering, writes files with Write tool, does not invoke `roadmapctl plan-paths` or `roadmapctl materialize`.
5. `rootline validate` and `roadmapctl check --strict` pass after a test materialization.

---

## Files Modified

| Repo | Path | Change |
|------|------|--------|
| rootline | `internal/rules/rules.go:160` | Add `NextByPattern` field |
| rootline | `internal/rules/describe.go:104-140` | Fix sort + `computeAllNextSequences` |
| rootline | `internal/rules/compute_sequence_test.go` | Tests: multi-pattern + determinism |
| rootline | `internal/rules/describe_test.go` | Test: `next_by_pattern` in JSON output |
| rootline | `.claude/skills/rootline/SKILL.md` | Note: multi-pattern `next` behavior |
| rootline | `.claude/skills/rootline/ref-schema.md` | Section: `next_by_pattern` |
| rootline | `docs/describe.md` | Example: `--field schema.id.next_by_pattern` |
| roadmapctl | `internal/cli/pathplan.go` | Delete |
| roadmapctl | `internal/cli/pathplan_test.go` | Delete |
| roadmapctl | `internal/cli/materialize.go` | Delete |
| roadmapctl | `internal/materialize/` (6 files) | Delete package |
| roadmapctl | `internal/roadmap/numbering.go` | Delete |
| roadmapctl | `internal/roadmap/numbering_plan_test.go` | Delete |
| roadmapctl | `docs/materialize-plan-schema.md` | Delete |
| roadmapctl | `internal/cli/cli.go:75-76` | Remove `newMaterializeCommand` + `newPlanPathsCommand` |
| roadmapctl | `docs/cli-contract.md` | Remove materialize section |
| roadmapctl | `internal/cli/cli_test.go` | Remove materialize tests |
| roadmapctl | `internal/cli/schema_compatibility_test.go` | Remove materialize refs |
| roadmapctl | `internal/roadmap/structure.go` | Remove materialize import/reference (check if import or comment; delete accordingly) |
| roadmapctl | `.claude/skills/roadmap/SKILL.md:70-80` | Update write invariant |
| roadmapctl | `.claude/skills/roadmap/plan-subcommand.md` | Full rewrite |
| roadmapctl | `.claude/skills/roadmap/common-logic.md` | Surgical edits |
| roadmapctl | `.claude/skills/roadmap/outcome-guide.md:16,42` | Update path-planning refs |
