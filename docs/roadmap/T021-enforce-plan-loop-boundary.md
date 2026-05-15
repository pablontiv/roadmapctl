---
estado: Specified
tipo: task
---
# T021: Enforce plan/loop boundary in plan-subcommand.md

**Contribuye a**: el subcomando `/roadmap plan` termina estrictamente en la creación de archivos `.md` — no ejecuta tasks ni inicia transiciones

## Preserva

- INV1: `roadmapctl check --strict` verde
  - Verificar: `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict`

## Contexto

Después de que `/roadmap plan` creó T019.md y T020.md, el agente continuó llamando
`roadmapctl transition start`, creando CLAUDE.md e implementando loop-subcommand.md
sin pasar por el loop. El STOP guard existe en Fase 4 pero no nombra explícitamente
las acciones prohibidas — el agente las interpretó como autorizadas.

El invariante es: la aprobación del árbol propuesto autoriza solo la creación de
archivos `.md`. Todo lo que sigue (transiciones, implementación de contenido, edición
de archivos descritos por las tasks) es responsabilidad exclusiva de `/roadmap loop`.

## Alcance

**In**:
1. Agregar inmediatamente después de `## Fase 4: Commit` en
   `/home/shared/roadmapctl/.claude/skills/roadmap/plan-subcommand.md`
   una sección de prohibición explícita con las acciones bloqueadas nombradas
   individualmente.

**Out**:
- No cambiar código
- No modificar otras secciones del skill

## Estado inicial esperado

`plan-subcommand.md § Fase 4` termina con "STOP. Informar: 'Archivos de planificación
creados...'" sin nombrar qué está prohibido hacer después.

## Criterios de Aceptación

- `grep -n "Prohibición\|transition start\|transition complete\|Prohibido" /home/shared/roadmapctl/.claude/skills/roadmap/plan-subcommand.md` retorna al menos una línea
- La sección nombra explícitamente: `transition start`, `transition complete`, modificar archivos descritos por las tasks, y continuar implementando
- `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict` verde

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/plan-subcommand.md` — archivo a modificar
