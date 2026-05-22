# Research — Scope mismatch en propuestas del skill `retrospective`

**Fecha:** 2026-05-21
**Status:** Research (no decisión todavía)
**Trigger:** "tenemos un problema con .claude/skills/retrospective/. actualmente todos los 'items' de la retro son para el repo en el que se ejecutan, lo cual no está mal, excepto cuando las propuestas son hacia los skills que se usan en múltiples proyectos"

---

## 1. Problema observado

El skill `retrospective` (definido en `.claude/skills/retrospective/SKILL.md`) produce en Fase 4 una **tabla de propuestas de corrección** cuyo formato actual asume que todos los artefactos están en el repo donde se ejecuta la retro:

```
| # | Tipo | Artefacto | Sección | Cambio propuesto | Previene |
```

La columna "Artefacto" recibe paths relativos como `.claude/skills/roadmap/loop-subcommand.md`. El skill no diferencia entre:

- **Artefactos locales** (del repo donde corre la retro): `CLAUDE.md`, `docs/roadmap/...`, skills propios del proyecto.
- **Artefactos upstream** (skills sincronizados desde otro repo): `roadmap`, `retrospective`, `integrate`, etc. — viven en `/home/shared/roadmapctl/.claude/skills/` y se distribuyen a `~/.claude/skills/` vía `scripts/sync-roadmap-skill.sh --install`.

Si la propuesta apunta a un skill upstream pero se materializa en el repo local, el cambio o no se aplica (archivo no existe en local) o se aplica a `~/.claude/skills/` y queda sobreescrito en el próximo sync.

## 2. Modelo de distribución de skills (estado actual)

Confirmado leyendo `scripts/sync-roadmap-skill.sh`:

- **Source-of-truth**: `/home/shared/roadmapctl/.claude/skills/<skill>/`
- **Destino del sync**: `~/.claude/skills/<skill>/` (user scope)
- **Modos del script**: `--install` (copia) y `--check` (verifica drift sin escribir)
- **Skills cubiertos**: `roadmap` por default; `--all` recorre todos los `SKILL.md` bajo `.claude/skills/`

El SKILL.md de `retrospective` ya declara su origen en frontmatter:
```yaml
source: pablontiv/roadmapctl
name: retrospective
```

Esto es el marcador declarativo natural para clasificar scope (upstream vs local).

Memory entry relacionada (`reference_roadmap_skill_source.md`):
> skill source-of-truth is this repo's `.claude/skills/roadmap/`, not `/opt/praxis`; sync via `scripts/sync-roadmap-skill.sh --install`

## 3. Evidencia histórica (vía backscroll)

### 3.1 Caso confirmatorio — picokit / Loop O01-picokit-module

Archivo: `~/.claude/projects/-home-shared-picokit/94f05d83-5aa4-4e91-8dfc-185bfdfb56cb.jsonl`

Tabla emitida en la retro:

| # | Tipo | Artefacto | Sección | Cambio | Previene |
|---|------|-----------|---------|--------|----------|
| 1 | AC | Tasks futuras Go en picokit | ACs | patrón `writeXxxWith` | E3 |
| 2 | Entorno | `CLAUDE.md` de picokit | Setup | `git remote add origin ...` | E1 |
| 3 | AC | `CLAUDE.md` de picokit | Convenciones | guardar exit code `go vet` | E2 |
| **4** | **Comprensión** | **`.claude/skills/roadmap/loop-subcommand.md`** | **Fase 3 — Implementar** | **agentes no deben declarar "untestable" paths OS-level Go sin intentar...** | **C1, C2** |

**Las propuestas 1–3 son locales a picokit; la #4 apunta a un skill upstream (roadmap) que vive en roadmapctl.** Si se materializa en picokit, no tiene efecto persistente. La corrección real debe pushearse a roadmapctl.

### 3.2 Caso mixto sin distinción visible — roadmapctl / Loop O22+O23

Archivo: `~/.claude/projects/-home-shared-roadmapctl/334403dd-af3e-4c7d-a5f8-b3e86dae80dc.jsonl`

| # | Artefacto | Naturaleza |
|---|-----------|------------|
| 1 | `.claude/skills/roadmap/loop-subcommand.md` | Upstream (pero ejecutado *en* upstream → válido localmente, casualmente) |
| 2 | `.claude/skills/roadmap/loop-subcommand.md` | Idem |
| 3 | `.claude/skills/roadmap/loop-subcommand.md` | Idem |
| 4 | `CLAUDE.md` (proyecto roadmapctl) | Local puro |

En este caso *funciona* porque roadmapctl es el repo upstream del skill `roadmap`, pero la tabla no expresa esa coincidencia. Si la misma retro corriera en otro proyecto, los items 1–3 fallarían silenciosamente.

### 3.3 Caso sin upstream-skill — backscroll / Loop T014-T018

Archivo: `~/.claude/projects/-home-shared-backscroll/f9a3bf3d-13da-434a-add1-cbdd34f2d1bb.jsonl`

Única propuesta apuntaba a `docs/roadmap/T015-update-readme-multi-source.md` — local puro, sin ambigüedad. Este es el caso donde el skill actual funciona bien.

## 4. Diagnóstico

El skill actual tiene la sección "Mapeo de artefactos por tipo de error" (líneas 142-152 de SKILL.md):

| Tipo de Error | Artefacto candidato |
|---------------|---------------------|
| Instrucción de skill incorrecta / ambigua | `.claude/skills/<skill>/SKILL.md` o subcomando |
| Regla operativa mal aplicada | `.claude/rules/<rule>.md` |
| ... |

Este mapeo **no distingue entre skills locales y skills upstream**. El path sugerido es relativo al repo donde corre la retro, lo cual es incorrecto cuando el skill se origina en otro repo.

Tampoco existe en Fase 3 ("Verificación pre-propuesta") un paso que clasifique scope antes de proponer.

## 5. Vectores de detección posibles

Tres mecanismos viables (en orden de simplicidad):

**A. Frontmatter `source:`** — leer el SKILL.md candidato, parsear frontmatter. Si `source:` apunta a un repo distinto al actual (`$(git config remote.origin.url)` o `$(basename $PWD)`), marcar como upstream. Cero overhead — el campo ya existe.

**B. Diff con `~/.claude/skills/<skill>/`** — si el skill local es idéntico byte-a-byte a la copia user-scope, asumir que es synced. Útil como fallback si frontmatter está ausente.

**C. Lista hardcoded de skills compartidos** — mantener en el SKILL.md una lista (`roadmap`, `retrospective`, `integrate`, `delegate`, etc.) con su repo upstream. Menos elegante pero cero ambigüedad.

## 6. Decisiones abiertas

1. **Formato de salida cuando hay propuestas upstream:**
   - (a) Dos tablas separadas ("Propuestas locales" / "Propuestas upstream")
   - (b) Tabla única con columnas `Scope` + `Repo destino`
   - (c) Solo emitir las locales; las upstream van como "observaciones" sin convertirlas en propuestas

2. **Comando de materialización para upstream:**
   - ¿Debería el skill emitir un handoff explícito tipo `cd /home/shared/roadmapctl && /roadmap plan ...`?
   - ¿O sólo señalar y dejar que el usuario decida (PR, edit directo, otro loop)?

3. **Detección del repo upstream a partir de `source: pablontiv/roadmapctl`:**
   - El skill necesitaría un mapping `pablontiv/roadmapctl → /home/shared/roadmapctl` (path local del clone)
   - ¿Asumir convención `~/shared/<repo-name>` o `~/shared/<basename>`? ¿Tener un registry en el skill?

4. **Alcance de la corrección:**
   - ¿Sólo skills, o también `.claude/rules/` (si llegaran a sincronizarse desde upstream en el futuro)?
   - ¿Aplicar misma lógica al skill `roadmap` (Fase plan), que también genera propuestas de artefactos?

## 7. Próximo paso

Pendiente decisión del usuario sobre §6 antes de pasar a fase de diseño/implementación.

## Referencias

- Skill actual: `.claude/skills/retrospective/SKILL.md`
- Script de sync: `scripts/sync-roadmap-skill.sh`
- Memory: `~/.claude/projects/-home-shared-roadmapctl/memory/reference_roadmap_skill_source.md`
- Retros analizadas:
  - picokit O01: `~/.claude/projects/-home-shared-picokit/94f05d83-5aa4-4e91-8dfc-185bfdfb56cb.jsonl`
  - roadmapctl O22+O23: `~/.claude/projects/-home-shared-roadmapctl/334403dd-af3e-4c7d-a5f8-b3e86dae80dc.jsonl`
  - backscroll T014-T018: `~/.claude/projects/-home-shared-backscroll/f9a3bf3d-13da-434a-add1-cbdd34f2d1bb.jsonl`
