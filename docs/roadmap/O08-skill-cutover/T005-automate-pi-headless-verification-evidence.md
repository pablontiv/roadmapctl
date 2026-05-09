---
estado: Pending
tipo: task
---
# T005: Automatizar evidencia Pi headless

**Outcome**: [O08 Cutover de skill](README.md)
**Contribuye a**: CE1, CE2

## Preserva

- INV1: No alcanza con grep para cambios de skill/guard.
  - Verificar: se ejecutan escenarios Pi headless.

## Contexto

`docs/roadmap-skill-integration.md` exige pruebas Pi headless para cambios al skill o guards. A medida que haya más comandos roadmapctl, conviene automatizar captura de evidencia.

## Alcance

**In**:
1. Crear script o procedimiento reproducible para escenarios headless.
2. Capturar comandos corridos y resultado de doctor/check/context/transition/materialize según aplique.
3. Documentar dónde guardar evidencia.
4. Integrar con checklist de release si es viable.

**Out**:
- Hacer Pi headless parte obligatoria de todos los tests unitarios.
- Modificar runtime Pi.

## Estado inicial esperado

- Docs contienen comandos manuales headless.

## Criterios de Aceptación

- Hay un comando/procedimiento único para correr escenarios.
- Evidencia demuestra no writes en escenarios preflight.
- Fallos bloquean release/cutover según checklist.

## Fuente de verdad

- `docs/roadmap-skill-integration.md`
- `.claude/skills/roadmap/SKILL.md`
- `scripts/sync-roadmap-skill.sh`
