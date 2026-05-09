---
estado: Specified
tipo: task
---
# T002: Make roadmap plan materialization token-light

**Outcome**: [Stabilize roadmap skill routing and materialization UX](README.md)

## Preserva

- roadmapctl remains the owner of canonical roadmap writes.
- Preflight, dry-run review, explicit approval, apply, and postcheck remain mandatory.

## Contexto

The current materialize flow can consume many tokens by forcing the LLM to serialize and read large JSON plans and dry-run reports with content and diffs.

## Alcance

**In**:
1. Revise plan-subcommand.md and common-logic.md to require temp plan/dry-run files.
2. Specify concise dry-run review fields: summary, diagnostics, path, operation, applied, and preconditions.
3. State that changes[].content and full diffs are read only on explicit request or targeted troubleshooting.

**Out**:
1. No removal of roadmapctl materialize.
2. No change to the materialize JSON schema unless explicitly approved later.

## Estado inicial esperado

plan-subcommand.md describes materialize dry-run/apply but does not strongly prevent large JSON/content/diff from entering the prompt context.

## Criterios de Aceptación

- Skill docs instruct agents to keep full plan and dry-run JSON in temp files by default.
- Normal dry-run reporting is limited to status, diagnostics, and planned paths.
- Headless materialize dry-run evidence shows no file modifications and no full content/diff dump in the final answer.

## Fuente de verdad

- .claude/skills/roadmap/plan-subcommand.md
- .claude/skills/roadmap/common-logic.md
- docs/materialize-plan-schema.md
- internal/cli/materialize.go
