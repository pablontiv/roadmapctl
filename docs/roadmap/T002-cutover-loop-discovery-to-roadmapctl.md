---
estado: Pending
tipo: task
---
# T002: Cut over loop discovery to roadmapctl pending/next

[[blocked_by:./T001-fix-outcome-stem-estado-requirement.md]]

## Preserva

- loop still functions correctly end-to-end

## Contexto

loop-subcommand.md was written before roadmapctl pending/next/decision existed. Those commands now exist and should be the canonical route for task queue discovery.

## Alcance

**In**:
1. loop-subcommand.md line 23: rootline tree -> roadmapctl pending --output json
2. loop-subcommand.md line 41: rootline graph -> roadmapctl next
3. loop-subcommand.md line 49: rootline query -> roadmapctl pending --output json
4. common-logic.md: mark rootline graph/query/tree as legacy/troubleshooting

**Out**:
1. no Go code changes
2. no changes to roadmapctl binary

## Estado inicial esperado

loop-subcommand.md uses rootline directly for discovery

## Criterios de Aceptación

- loop-subcommand.md has no rootline graph, rootline query, rootline tree in main flow
- roadmapctl pending and roadmapctl next used for workspace selection and queue
- common-logic.md marks those rootline commands as legacy/troubleshooting only

## Fuente de verdad

- .claude/skills/roadmap/loop-subcommand.md
- .claude/skills/roadmap/common-logic.md
