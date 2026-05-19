---
estado: Specified
tipo: task
---
# T002: Sweep skill markdown for `--roadmap-root` mentions

**Outcome**: [O22 Drop --roadmap-root flag and document workspace](README.md)
**Contribuye a**: skill `/roadmap` deja de pasar el flag eliminado por T001

[[blocked_by:./T001-drop-roadmap-root-flag-from-go.md]]

## Preserva

- INV1: el skill `/roadmap` sigue funcionando contra el roadmap del propio repo
  - Verificar: tras los cambios, `roadmapctl bootstrap --repo . --output json` desde cualquier ejemplo del skill no produce error "unknown flag"
- INV2: la lógica del skill no cambia — solo se eliminan menciones del flag y se ajustan los comandos de ejemplo

## Contexto

El skill `/roadmap` vive en `.claude/skills/roadmap/` y consiste en múltiples archivos markdown: `SKILL.md`, `loop-subcommand.md`, `plan-subcommand.md`, `pending-subcommand.md`, `decision-tree-subcommand.md`, `autonomous-mode.md`, `framework-reference.md`, `common-logic.md`, `outcome-guide.md`, `task-guide.md`, `pr-workflow.md`.

Estos archivos contienen ejemplos de comandos como:

```bash
roadmapctl bootstrap --repo <repo-path> --roadmap-root <roadmap-root-si-se-conoce> --output json
roadmapctl doctor --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl transition can-start <path> --repo <repo> --roadmap-root docs/roadmap
roadmapctl transition start <path> --repo <repo> --roadmap-root docs/roadmap --apply
roadmapctl transition complete <path> --repo <repo> --roadmap-root docs/roadmap --apply
```

Tras T001 estos comandos fallan con "unknown flag: --roadmap-root". Hay que eliminar el flag de todos los ejemplos. El `roadmap-root` queda implícito (hardcodeado en `docs/roadmap` en Go).

Buscar todas las menciones:

```bash
grep -rn "roadmap-root" /home/shared/roadmapctl/.claude/skills/roadmap/
```

Reemplazar removiendo el flag y el placeholder/valor que lo acompaña. Ejemplo:

- Antes: `roadmapctl bootstrap --repo <repo-path> --roadmap-root <roadmap-root-si-se-conoce> --output json`
- Después: `roadmapctl bootstrap --repo <repo-path> --output json`

Si una sección del skill habla del concepto roadmap-root (no del flag), reformular para que refleje la convención: `docs/roadmap/` es la ubicación fija; viene resuelta por `roadmapctl bootstrap` en el campo `roadmap_root` del JSON.

## Alcance

**In**:
1. `grep -rn "roadmap-root" .claude/skills/roadmap/` para enumerar todas las menciones
2. Por cada match, decidir:
   - Si es invocación CLI con `--roadmap-root <val>`: eliminar el flag y su valor
   - Si es texto explicativo sobre el concepto "roadmap-root": reformular como "convención fija `docs/roadmap/`" o "valor del campo `roadmap_root` del JSON de bootstrap"
   - Si es placeholder en tabla (e.g. `<roadmap-root>`): mantener si refiere al valor resuelto por bootstrap; eliminar si refiere al flag
3. Verificar que después del sweep, el contenido de los archivos sigue siendo coherente y los comandos de ejemplo son válidos contra la versión post-T001 del binario

**Out**:
- Cambios al README.md o docs/cli-contract.md — T003/T004
- Cambios estructurales al skill (lógica de routing, fases, etc.) — fuera del scope del outcome
- Sincronización del skill al repo fuente vía sync-roadmap-skill.sh — debe ejecutarse pero no es un cambio editorial

## Estado inicial esperado

- `grep -rn "roadmap-root" .claude/skills/roadmap/ | wc -l` retorna un número positivo (varias menciones existen)
- T001 ya está Completed: `roadmapctl --help` no muestra el flag, pasarlo da "unknown flag"

## Criterios de Aceptación

- `grep -rn "\-\-roadmap-root" .claude/skills/roadmap/` retorna 0 matches (cero invocaciones del flag)
- `grep -rn "roadmap-root" .claude/skills/roadmap/` solo retorna menciones que refieren al concepto/campo `roadmap_root` (no al flag CLI) o están en bloques de YAML/TOML legacy (e.g. `roadmap-root: docs/roadmap` como frontmatter legacy citado como referencia histórica)
- Verificación headless con Pi (per protocolo del skill): los dos prompts de verificación documentados en SKILL.md sección "Verificación obligatoria al modificar este skill" ejecutan sin errores `unknown flag`
- `./scripts/sync-roadmap-skill.sh --install` ejecuta sin errores

## Fuente de verdad

- `.claude/skills/roadmap/SKILL.md`
- `.claude/skills/roadmap/loop-subcommand.md`
- `.claude/skills/roadmap/plan-subcommand.md`
- `.claude/skills/roadmap/pending-subcommand.md`
- `.claude/skills/roadmap/decision-tree-subcommand.md`
- `.claude/skills/roadmap/autonomous-mode.md`
- `.claude/skills/roadmap/framework-reference.md`
- `.claude/skills/roadmap/common-logic.md`
- `.claude/skills/roadmap/outcome-guide.md`
- `.claude/skills/roadmap/task-guide.md`
- `.claude/skills/roadmap/pr-workflow.md`
- `scripts/sync-roadmap-skill.sh` (ejecución, no edición)
