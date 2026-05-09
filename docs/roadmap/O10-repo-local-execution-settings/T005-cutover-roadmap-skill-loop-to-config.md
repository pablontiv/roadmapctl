---
estado: Completed
tipo: task
---
# T005: Cut over roadmap skill loop to config-driven behavior

**Outcome**: [Repo-local roadmap execution settings](README.md)

[[blocked_by:./T002-expose-execution-settings-in-context.md]]
[[blocked_by:./T004-add-roadmap-context-compaction-extension.md]]
[[blocked_by:./T007-implement-parallel-batch-materialization.md]]

## Preserva

- roadmapctl doctor/check remain mandatory before loop execution.
- The skill still uses roadmapctl next as the canonical ready/blocked source.
- The skill does not parse legacy config or migrate config itself.

## Contexto

The user wants behavior settings fixed locally per repo. The skill must stop advertising behavior flags and must use context JSON values instead.

## Alcance

**In**:
1. Update .claude/skills/roadmap/SKILL.md bootstrap/config sections to treat legacy as migration input handled by roadmapctl, not fallback for implemented flows.
2. Update .claude/skills/roadmap/loop-subcommand.md to keep only --filter and --max, and to define config-driven loop_max_tasks, parallel, autonomy, compact_after_task_commit, and pr_mode behavior.
3. Update .claude/skills/roadmap/pr-workflow.md so PR mode comes from pr_mode config and merge behavior follows autonomy plus pr_merge_strategy.
4. Document opportunistic parallel waves using blocked_by/roadmapctl next as the only dependency source.
5. Document compact_roadmap_context with /compact fallback after durable commit/push/PR bookkeeping.
6. Update materialization-related skill language so roadmapctl-owned batch apply is allowed while direct multi-file skill writes remain forbidden.

**Out**:
1. No Go config parser changes in this task.
2. No extension implementation in this task.

## Estado inicial esperado

The skill currently documents --parallel, --worktree, --self-pace, --skip-reviews, --checkpoint-interval, and --pr as loop options. It says legacy remains a fallback while migration completes.

## Especificación Técnica

In bootstrap, list the new context JSON operational fields the skill must read. In loop discovery, compute effective_max as --max when present, else loop_max_tasks, with 0 meaning unlimited. Remove the confirmation step for supervised/until_done. Keep manual confirmation for manual mode. State that structural blocked_by repair must be followed by roadmapctl check --strict before continuing; if no deterministic safe edit path is available, even until_done must stop and report the required dependency.

## Criterios de Aceptación

- rg -- '--parallel|--worktree|--self-pace|--skip-reviews|--checkpoint-interval|--pr' .claude/skills/roadmap returns no active loop option documentation except historical references explicitly marked obsolete if any are retained.
- .claude/skills/roadmap/loop-subcommand.md says /roadmap loop accepts only --filter and --max plus global --repo.
- Loop docs define autonomy values manual, supervised, and until_done exactly as in the design spec.
- Loop docs define parallel waves based only on roadmapctl next/blocked_by and explain conflict repair behavior by autonomy mode.
- Loop docs define compact_after_task_commit order and fallback to /compact when compact_roadmap_context is unavailable.
- Plan/materialization docs allow a single roadmapctl-owned batch apply command when roadmapctl guarantees canonical writes, per-file diagnostics, and postcheck.

## Fuente de verdad

- docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md
- .claude/skills/roadmap/SKILL.md
- .claude/skills/roadmap/loop-subcommand.md
- .claude/skills/roadmap/pr-workflow.md
- .claude/skills/roadmap/common-logic.md
