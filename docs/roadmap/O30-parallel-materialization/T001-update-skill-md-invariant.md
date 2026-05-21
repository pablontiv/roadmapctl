---
estado: Specified
tipo: task
---
# T001: Actualizar invariante de escritura en `SKILL.md`

**Outcome**: [O30 Materialización paralela](README.md)
**Contribuye a**: Eliminar la restricción singleton-writer que impide despachar Agents en la fase de materialización

## Preserva

- INV1: El skill sigue siendo el coordinador — ningún Agent escribe sin aprobación explícita del usuario y preflight exitoso
  - Verificar: el texto nuevo mantiene "después de aprobación explícita y preflight pasado" como precondición

## Contexto

`SKILL.md` en `.claude/skills/roadmap/` declara actualmente:

> "El skill es el único writer vía Write tool, después de aprobación explícita y preflight pasado."

Esta restricción impide que el skill despache Agents para escritura paralela. La intención original era evitar escrituras no coordinadas (heredocs en shell, loops con `rootline new`), pero la redacción prohíbe el caso válido de Agents coordinados post-aprobación.

El cambio reemplaza "único writer" por "coordinador": el skill orquesta la materialización, los Agents ejecutan las escrituras.

Sección a modificar: `## Invariante de escritura segura` en `/home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md`.

## Alcance

**In**:
1. Reemplazar el párrafo "El skill es el único writer vía Write tool..." por el modelo coordinador
2. Mantener las prohibiciones existentes (heredocs, `rootline new`, escribir sin aprobación, escribir con preflight non-zero)
3. Agregar: "Permitido: despachar Agents coordinados para escritura paralela tras aprobación y preflight exitoso"

**Out**:
- No modificar ningún otro archivo en esta task
- No cambiar la lógica de `plan-subcommand.md` (eso es T002)

## Estado inicial esperado

- `grep "único writer" /home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` retorna al menos 1 match

## Criterios de Aceptación

- `grep "único writer" /home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` retorna 0 matches
- `grep "coordinador" /home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` retorna al menos 1 match
- `grep "aprobación explícita" /home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` retorna al menos 1 match (precondición preservada)

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` — sección `## Invariante de escritura segura`
