---
estado: Completed
tipo: task
---
# T010: Documentar integración obligatoria con el skill roadmap

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE3

[[blocked_by:./T001-define-cli-contract.md]]
[[blocked_by:./T006-implement-doctor-command.md]]
[[blocked_by:./T008-implement-rootline-backed-checks.md]]

## Preserva

- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: docs indican bloqueo, no auto-fix.

## Contexto

El usuario decidió que `roadmapctl` no es opcional para comandos implementados de `/roadmap`. Si un comando escribe, muta, ejecuta o declara validez del roadmap, debe pasar por `roadmapctl`.

## Alcance

**In**:
1. Documentar regla obligatoria: `doctor` y `check` antes de writes/ejecución.
2. Documentar postcheck obligatorio para comandos que materializan.
3. Documentar que modo conceptual puede existir sin `roadmapctl` solo si no escribe ni declara materialización.
4. Preparar snippet para actualizar `/roadmap plan` y `/roadmap loop`.
5. Documentar errores esperados cuando falta `roadmapctl` o `rootline`.

**Out**:
- Editar skills instalados en esta task.
- Implementar hooks de Pi o Claude.
- Cambiar Rootline.

## Estado inicial esperado

- Contrato CLI definido.
- `doctor`/`check` tienen comportamiento claro.

## Criterios de Aceptación

- Docs declaran que `roadmapctl` es requerido desde día 1 para comandos implementados que escriben, mutan o ejecutan.
- Docs incluyen comandos exactos de preflight y postcheck.
- Docs incluyen política de bloqueo si `roadmapctl` falla.
- La integración no sugiere fallback a markdown libre.

## Fuente de verdad

- `docs/roadmap-skill-integration.md`
- `.claude/skills/roadmap/SKILL.md` en el repo consumidor cuando se aplique
- `docs/cli-contract.md`
