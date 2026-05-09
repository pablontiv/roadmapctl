---
estado: Completed
tipo: task
---
# T004: Add roadmap context compaction extension

**Outcome**: [Repo-local roadmap execution settings](README.md)

## Preserva

- Compaction remains Pi runtime behavior, not roadmapctl behavior.
- The extension must be safe in non-interactive contexts and report errors clearly.

## Contexto

The design prefers a dedicated compact_roadmap_context tool and falls back to /compact instructions when the tool is unavailable.

## Alcance

**In**:
1. Create a project-local Pi extension under .pi/extensions/roadmap-context/index.ts or another documented project-local extension path.
2. Register a compact_roadmap_context tool that calls ctx.compact with roadmap-specific custom instructions.
3. Document how the extension is loaded/synced with the roadmap skill.

**Out**:
1. No automatic compaction trigger outside explicit tool invocation.
2. No roadmapctl code changes in this task.

## Estado inicial esperado

No project-local Pi extension exists for roadmap context compaction. Pi supports ctx.compact from extensions and /compact as an interactive fallback.

## Especificación Técnica

Use Type.Object from typebox for parameters. Accept optional task_path, commit_hash, validation_summary, next_work, and config_summary strings so the skill can pass iteration context. Build customInstructions deterministically from those fields plus a fixed roadmap continuation checklist. Call ctx.compact({ customInstructions, onComplete, onError }). Do not await compaction because ctx.compact is fire-and-forget.

## Criterios de Aceptación

- The extension TypeScript file registers a tool named compact_roadmap_context.
- The tool passes custom instructions preserving current roadmap goal, completed task path, commit hash, validation results, next task/wave state, unresolved blockers/conflicts, and relevant config values.
- The tool reports compaction queued/failed status without throwing on normal compaction callbacks.
- The roadmap skill docs can reference compact_roadmap_context by exact name.

## Fuente de verdad

- docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md
- /home/pones/.local/lib/node_modules/@earendil-works/pi-coding-agent/docs/extensions.md
- /home/pones/.local/lib/node_modules/@earendil-works/pi-coding-agent/docs/compaction.md
- /home/pones/.local/lib/node_modules/@earendil-works/pi-coding-agent/examples/extensions/trigger-compact.ts
