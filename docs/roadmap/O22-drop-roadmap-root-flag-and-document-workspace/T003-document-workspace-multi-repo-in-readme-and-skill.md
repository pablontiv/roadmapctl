---
estado: Specified
tipo: task
---
# T003: Document workspace multi-repo in README and SKILL.md

**Outcome**: [O22 Drop --roadmap-root flag and document workspace](README.md)
**Contribuye a**: usuarios entienden que cada repo del workspace mantiene su propio roadmap, cerrando los issues #2 y #3 por documentación

## Preserva

- INV1: documentación existente (Quick Start, Core Idea, etc.) no se altera ni reordena más allá de añadir una sección nueva en posición lógica
- INV2: SKILL.md Paso 0 sigue distinguiendo workspace vs single-repo via `test -d .git`

## Contexto

El modelo workspace existe en el código desde hace tiempo (`internal/cli/pending.go:114` define `workspaceRepoRoots()` que walks `.git` siblings). Está implementado y funcional. Lo que falta es **documentación explícita** del modelo de uso para evitar confusiones como las que originaron los GitHub issues #2 (esperar routing cross-repo desde un roadmap central) y #3 (esperar `.roadmapctl.toml` en repo root porque "code repos no tienen roadmap").

Investigación con backscroll y diálogo con usuario establecieron la convención: **cada repo participante en el workspace mantiene su propio roadmap completo bajo `<repo>/docs/roadmap/` con `.stem` + `.roadmapctl.toml` + outcomes + tasks**. No existe el escenario "code repo sin roadmap propio".

Ubicación de los cambios:

### `/home/shared/roadmapctl/README.md`

Nueva sección titulada **Workspace mode** insertada entre **Quick Start** (línea 36) y **Core Idea** (línea 60). Contenido sugerido:

```markdown
## Workspace mode

When `roadmapctl` runs from a directory without a `.git` directory but containing sibling repos with their own `.git`, it operates in **workspace mode**.

**Convention**: Each participating repo maintains its own complete roadmap under `<repo>/docs/roadmap/` with `.stem`, `.roadmapctl.toml`, outcomes, and tasks. Each repo is autonomous — its roadmap, its tasks, its commits.

**Layout example**:

```text
my-workspace/                    # parent dir without .git
├── docs/                        # repo 1: docs/roadmap/ + .stem + .roadmapctl.toml
│   ├── .git/
│   └── docs/roadmap/
│       ├── .stem
│       └── .roadmapctl.toml
├── tsg-valuecreation-core/      # repo 2: same layout
│   ├── .git/
│   └── docs/roadmap/
│       └── ...
└── tsg-valuecreation-frontend/  # repo 3: same layout
    ├── .git/
    └── docs/roadmap/
        └── ...
```

**Invocation**: Most commands operate on a single repo at a time. Pass `--repo <path>` to target a specific repo. `roadmapctl pending --workspace` iterates the discovered repos and aggregates results.

**Anti-pattern**: Do not create "code repos" that lack their own roadmap, expecting a central roadmap repo to commit changes on their behalf. Each repo is autonomous; cross-repo commit routing is not supported. If a repo needs to participate in the workspace, it needs its own `<repo>/docs/roadmap/`.
```

### `/home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md`

Extender la sección **Workspace mode** existente (alrededor de las líneas que dicen "Escanear subdirectorios inmediatos con `.git` + config roadmap") con clarificación operativa:

- En workspace mode, el skill se invoca **por repo** (`/roadmap loop` debe ejecutarse en cada repo, o usar `--repo <name>` para targetear uno específico)
- Cada repo gestiona su propio commit/push según el `auto_push` resuelto en su propio `.roadmapctl.toml`
- No existe routing de commits cross-repo: cada task de un repo debe tocar solo archivos de ese repo

Esta extensión refuerza la convención que el README documenta, desde la perspectiva del skill.

## Alcance

**In**:
1. Añadir sección "Workspace mode" al `README.md` entre Quick Start y Core Idea
2. Actualizar tabla de contenidos (`README.md` línea ~22-32) para incluir el nuevo anchor `#workspace-mode`
3. Extender SKILL.md Paso 0 / sección Workspace mode con la clarificación operativa
4. Verificar que el ejemplo de layout en README usa rutas reales/realistas y no contradice ningún otro doc

**Out**:
- Cambios al CLI o comportamiento Go — no aplica
- Cambios a `docs/cli-contract.md` — T004
- Eliminación del flag `--roadmap-root` — T001/T002
- Crear plantillas o `roadmapctl init` para auto-crear roadmaps en code repos — fuera del scope del outcome

## Estado inicial esperado

- `README.md` no contiene sección "Workspace mode"
- `grep -n "Workspace" .claude/skills/roadmap/SKILL.md` muestra menciones pero no clarifica la invocación per-repo ni el commit per-repo

## Criterios de Aceptación

- `README.md` contiene una sección con encabezado `## Workspace mode` que incluye: cuándo aplica, convención (cada repo su propio roadmap), bloque de layout de ejemplo, advertencia de anti-patrón
- Tabla de contenidos del README actualizada con link a `#workspace-mode`
- `.claude/skills/roadmap/SKILL.md` en la sección Workspace mode aclara explícitamente: (a) el loop se invoca por repo, (b) cada repo gestiona su propio commit/push según su `auto_push`
- Ambos archivos (README + SKILL.md) refuerzan el mismo mensaje: cada repo del workspace tiene su propio `docs/roadmap/`
- Ninguno menciona `code_repos = [...]` ni puntero workspace_roadmap ni modo minimal ni discriminación por `.stem`

## Fuente de verdad

- `README.md` — sección nueva + TOC
- `.claude/skills/roadmap/SKILL.md` — extensión de sección Workspace mode existente
