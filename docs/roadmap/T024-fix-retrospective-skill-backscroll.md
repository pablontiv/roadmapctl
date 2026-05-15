---
estado: Completed
tipo: task
---
# T024: Fix retrospective SKILL.md to use backscroll as primary retrieval

**Contribuye a**: la retrospectiva usa el binario `backscroll` como herramienta primaria de búsqueda en historial, no parseo manual de `.jsonl`

## Preserva

- INV1: `roadmapctl check --strict` verde
  - Verificar: `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict`

## Contexto

Al ejecutar la retrospectiva con argumento "usa backscroll", el agente usó el bloque
`PROJECT_SLUG` + python3 de "Recuperación de errores pre-compact" en vez del binario
`backscroll search`. Esto produjo:

1. Error en la construcción del PROJECT_SLUG (usó `echo /home/shared` en vez de `pwd`)
2. Scripts Python complejos en vez de un simple `backscroll search <keywords>`
3. Resultados incompletos — el usuario tuvo que señalar que faltaban propuestas

El binario `backscroll` existe, está instalado y está diseñado exactamente para
recuperar historial de sesiones. El bloque manual es un fallback de cuando backscroll
no está disponible, no la ruta primaria.

## Alcance

**In**:
1. Modificar `## Fase 0 — Recuperación de contexto`, subsección "Recuperación de
   errores pre-compact (si sesión larga)" en
   `/home/shared/roadmapctl/.claude/skills/retrospective/SKILL.md`
   para que use `backscroll search` como ruta primaria y el bloque `PROJECT_SLUG`
   solo como fallback cuando `backscroll` no está disponible.

**Out**:
- No eliminar el bloque PROJECT_SLUG — marcarlo como fallback
- No cambiar otras fases del skill

## Estado inicial esperado

La subsección muestra el bloque `PROJECT_SLUG=$(pwd | tr '/' '-'...)` + python3 como
la única ruta de recuperación, sin mención de `backscroll search`.

## Criterios de Aceptación

- `grep -n "backscroll search\|Primario\|Fallback\|fallback" /home/shared/roadmapctl/.claude/skills/retrospective/SKILL.md` retorna al menos 2 líneas
- El bloque `backscroll search` aparece ANTES del bloque `PROJECT_SLUG`
- El bloque `PROJECT_SLUG` tiene etiqueta explícita de fallback
- `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict` verde

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/retrospective/SKILL.md` — archivo a modificar
