---
estado: Completed
tipo: task
---
# T012: Versionar el skill roadmap y sincronizarlo al user scope

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE3 y CE4

[[blocked_by:./T001-define-cli-contract.md]]

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: el skill `roadmap` vive en este repo, no como lógica dentro de `rootline`.
- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: el hook solo sincroniza skill files; no ejecuta materialización.

## Contexto

La fuente canónica del skill `/roadmap` debe vivir en `/home/shared/roadmapctl/.claude/skills/roadmap/`. El user scope `~/.claude/skills/roadmap` debe ser una copia instalada por hook para que las sesiones usen la versión del repo.

## Alcance

**In**:
1. Versionar `.claude/skills/roadmap/` dentro del repo `roadmapctl`.
2. Añadir hook de git que copie `.claude/skills/roadmap` a `~/.claude/skills/roadmap`.
3. Hacer la sincronización idempotente y explícita.
4. Documentar que el repo `roadmapctl` es la fuente canónica del skill.
5. Verificar que la copia instalada coincide con la fuente del repo.

**Out**:
- Borrar manualmente skills de user scope sin backup.
- Cambiar el motor Rootline.
- Instalar binarios del MVP.

## Estado inicial esperado

- El repo contiene `.claude/skills/roadmap/`.
- El user scope puede contener una copia previa del skill.

## Criterios de Aceptación

- `.claude/skills/roadmap/SKILL.md` existe en el repo.
- El hook copia todos los archivos del skill a `~/.claude/skills/roadmap`.
- El hook no toca otros skills.
- Un comando de verificación confirma que fuente e instalación coinciden.
- README o docs indican que `roadmapctl` es el hogar canónico del skill `roadmap`.

## Fuente de verdad

- `.claude/skills/roadmap/`
- `.githooks/post-merge` o hook equivalente
- `README.md`
- `docs/roadmap/O01-roadmapctl-mvp/README.md`
