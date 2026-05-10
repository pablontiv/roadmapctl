# Materialize plan input schema

`roadmapctl materialize` is planned as a deterministic writer for already-approved roadmap plans. It does **not** decompose free-form chat, ask an LLM for task boundaries, or infer missing requirements. The `/roadmap plan` skill remains responsible for AI-assisted decomposition and user approval, then passes structured JSON to roadmapctl.

## Versioning

Every input document is versioned:

```json
{
  "version": 1,
  "kind": "roadmapctl/materialize-plan",
  "items": []
}
```

Rules:

- `version` is required and must be `1` for the initial contract.
- `kind` is required and must be `roadmapctl/materialize-plan`.
- Unknown top-level fields are ignored for forward compatibility unless they conflict with known fields.
- Future incompatible changes must use a new integer `version`.

## Top-level shape

```json
{
  "version": 1,
  "kind": "roadmapctl/materialize-plan",
  "title": "Optional human summary",
  "items": [
    {
      "type": "outcome",
      "slug": "config-context-workspace",
      "title": "Config/context/workspace",
      "description": "Outcome context for generated README.",
      "acceptance_criteria": ["Context command reports effective config."],
      "tasks": [
        {
          "slug": "implement-context-command",
          "title": "Implement context command",
          "description": "Expose effective roadmapctl context as JSON.",
          "preserves": ["Command is read-only."],
          "context": "The skill bootstrap needs stable config data.",
          "scope_in": ["Add roadmapctl context."],
          "scope_out": ["No mutation."],
          "initial_state": "Config loader exists.",
          "acceptance_criteria": ["JSON has kind roadmapctl/context."],
          "source_of_truth": ["internal/cli/context.go"],
          "blocked_by": []
        }
      ]
    },
    {
      "type": "task",
      "slug": "update-release-docs",
      "title": "Update release docs",
      "description": "Document release checklist.",
      "preserves": ["Docs-only change."],
      "context": "Release governance needs explicit steps.",
      "scope_in": ["Update docs/release.md."],
      "scope_out": ["No Go code."],
      "initial_state": "Release doc exists.",
      "acceptance_criteria": ["Checklist includes coverage gate."],
      "source_of_truth": ["docs/release.md"],
      "blocked_by": []
    }
  ]
}
```

`items[]` may contain outcomes and direct tasks. Outcome tasks live in the outcome's `tasks[]` array. The materializer must not create a single markdown summary file for multiple tasks.

## Outcome object

| Field | Required | Type | Meaning |
|-------|----------|------|---------|
| `type` | yes | string | Must be `outcome`. |
| `slug` | yes | string | Human slug without numeric prefix. Lowercase kebab-case recommended. |
| `title` | yes | string | Outcome title for README heading. |
| `description` | yes | string | Context paragraph for README. |
| `acceptance_criteria` | yes | array(string) | Observable outcome-level criteria. |
| `tasks` | yes | array(task) | One or more child tasks. |
| `contributes_to` | no | array(string) | Outcome criterion labels if supplied by the plan. |

Validation:

- `tasks` must contain at least one task.
- `slug` must be portable and must not include path separators or numeric `OXX-` prefix.
- The materializer assigns `OXX` numbers deterministically from the current roadmap root.

## Task object

Task objects are used both inside outcomes and directly under `items[]` with `type: "task"`.

| Field | Required | Type | Meaning |
|-------|----------|------|---------|
| `type` | direct only | string | Direct task must set `type: "task"`; nested outcome tasks may omit it or set `task`. |
| `slug` | yes | string | Human slug without numeric `TXXX-` prefix. |
| `title` | yes | string | Task title for heading. |
| `description` | yes | string | Short implementation description. |
| `preserves` | yes | array(string) | Invariants and preservation checks. |
| `context` | yes | string | Why the task exists. |
| `scope_in` | yes | array(string) | Included work. |
| `scope_out` | yes | array(string) | Explicit exclusions. |
| `initial_state` | yes | string | Expected state before execution. |
| `acceptance_criteria` | yes | array(string) | Observable ACs; at least one required. |
| `source_of_truth` | yes | array(string) | Paths/docs governing implementation; at least one required. |
| `blocked_by` | no | array(dependency) | Dependencies expressed against plan-local refs or existing paths. |
| `technical_spec` | no | string | Optional deterministic details to include as `## Especificación Técnica`. |

Validation:

- Required string fields must be non-empty after trimming whitespace.
- Required arrays must be non-empty and contain non-empty strings.
- `slug` must be portable and must not include path separators or numeric `TXXX-` prefix.
- The materializer assigns `TXXX` numbers deterministically within the destination scope.

## Dependencies

Dependencies use explicit objects so roadmapctl can resolve them after numbering:

```json
{
  "ref": "config-context-workspace/implement-config-loader",
  "path": "../O03-config-context-workspace/T002-implement-config-discovery-and-toml-loader.md"
}
```

| Field | Required | Meaning |
|-------|----------|---------|
| `ref` | conditional | Plan-local reference, usually `<outcome-slug>/<task-slug>` or `<task-slug>` for direct tasks. |
| `path` | conditional | Existing roadmap task path relative to the task being written or roadmap root. |

Exactly one of `ref` or `path` is required.

Rules:

- `ref` must resolve to another task in the same materialize plan.
- `path` is validated during dry-run and must point to an existing or concurrently planned `TXXX-*.md` task.
- Emitted markdown must use explicit relative `[[blocked_by:...]]` links.
- Bare basename dependencies are invalid.

## Diagnostics

Planned input validation diagnostics:

| ID | Severity | Meaning |
|----|----------|---------|
| `RMC_MATERIALIZE_INPUT_VERSION_UNSUPPORTED` | error | `version` is missing or unsupported. |
| `RMC_MATERIALIZE_INPUT_KIND_INVALID` | error | `kind` is missing or not `roadmapctl/materialize-plan`. |
| `RMC_MATERIALIZE_INPUT_EMPTY` | error | No outcomes or tasks were provided. |
| `RMC_MATERIALIZE_INPUT_FIELD_MISSING` | error | Required field is absent or empty. |
| `RMC_MATERIALIZE_INPUT_SLUG_INVALID` | error | Slug is not portable or includes path/numeric prefix. |
| `RMC_MATERIALIZE_INPUT_DEPENDENCY_INVALID` | error | Dependency object has neither/both `ref` and `path`, or uses a bare target. |
| `RMC_MATERIALIZE_INPUT_DEPENDENCY_UNRESOLVED` | error | Dependency cannot be resolved to a plan-local, concurrently planned, or existing task. |
| `RMC_MATERIALIZE_PLAN_CONFLICT` | error | Planned path collides with an existing unrelated roadmap item. |

Additional target-apply diagnostics:

| ID | Severity | Meaning |
|----|----------|---------|
| `RMC_MATERIALIZE_TARGET_UNKNOWN` | error | `--target` is absent from the frozen change-set. |
| `RMC_MATERIALIZE_TARGET_DUPLICATE` | error | `--target` appears multiple times in the frozen change-set. |
| `RMC_MATERIALIZE_TARGET_INVALID` | error | `--target` is not a single canonical roadmap markdown create change. |

All diagnostics use the standard report shape with `kind: roadmapctl/materialize` once the command exists. `details` should include JSON pointer-like paths such as `/items/0/tasks/1/slug` when possible.

## Frozen change-set apply

`roadmapctl materialize --plan <plan-json> --dry-run --output json` returns `changes[]` with deterministic `path`, `operation`, `content`, `diff`, and preconditions. Saving that JSON creates a frozen change-set that can be applied in one roadmapctl-owned batch:

```bash
roadmapctl materialize --changes dry-run.json --apply --repo <repo> --roadmap-root <roadmap-root> --output json
```

Rules:

- `--changes` requires `--apply`; `--target` is optional.
- Without `--target`, roadmapctl validates the whole frozen change-set before writing, orders parent/container changes before child tasks, writes all allowed changes, validates created markdown, and runs the standard postcheck before reporting success.
- Batch apply accepts only allowlisted bootstrap changes (`.`, `.stem`, `.roadmapctl.toml`) and canonical roadmap markdown files (`TXXX-*.md`, `OXX-*/README.md`, or `OXX-*/TXXX-*.md`).
- Existing planned paths fail before writing via `RMC_MATERIALIZE_PLAN_CONFLICT`, and diagnostics identify the concrete blocking path.
- With `--target`, the target must match exactly one `changes[].path` entry; target apply writes only that selected canonical markdown file, does not recompute numbering from the plan, and does not create sibling roadmap files.

## Skill integration

This document is the canonical schema source inside the `roadmapctl` repository. Installed roadmap skills and consuming repositories must not assume this file exists relative to their own working tree. A future `roadmapctl materialize schema --output json` command may expose a machine-readable schema, but that CLI surface is deferred until explicitly approved.

The `/roadmap plan` skill must:

1. decompose and present the proposed tree to the user;
2. stop for explicit approval;
3. serialize the approved tree to this JSON shape in a temporary plan file rather than pasting large JSON into chat;
4. pass the temp plan file to roadmapctl for dry-run materialization and save stdout to a temporary dry-run JSON file;
5. review dry-run output normally by extracting only `summary`, `diagnostics`, `path`, `operation`, `applied`, and `preconditions`; inspect `changes[].content` or full diffs only on explicit user request or targeted troubleshooting;
6. save the dry-run JSON when using a frozen change-set and apply approved files with roadmapctl-owned batch apply (`--plan ... --apply` or `--changes <dry-run-json> --apply`);
7. avoid writing roadmap markdown directly once roadmapctl materialization exists. Granular `--target` apply is reserved for recovery/troubleshooting or explicit one-file approval.

roadmapctl must reject free-form prose input. If the skill lacks enough information to populate required fields, the skill asks the user before calling roadmapctl.
