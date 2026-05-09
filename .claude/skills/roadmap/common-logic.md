# Lógica Común — Materialización y Ejecución

> En workspace mode, `<roadmap-root>` se reemplaza por `<abs-roadmap-root>` y los comandos git usan `git -C <repo-path>`.

## Guard obligatorio: roadmapctl

Para cualquier flujo que escriba, mute, ejecute tasks o declare validez del roadmap, `roadmapctl` es obligatorio además de Rootline.

Antes de escribir, mutar o ejecutar:

```bash
command -v roadmapctl
roadmapctl doctor --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Después de cualquier materialización o mutación del roadmap:

```bash
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si `roadmapctl` no existe o cualquier comando sale non-zero, detenerse antes de continuar. Reportar comando, exit code y diagnostic IDs si hubo JSON. No auto-fix, no fallback a markdown libre, no ejecutar tasks y no commitear mutaciones del roadmap.

La planificación conceptual que no escribe, no muta, no ejecuta y no declara validez puede continuar sin `roadmapctl`.

## Modelo

El roadmap usa máximo dos niveles:

```text
<roadmap-root>/
├── O01-outcome/
│   ├── README.md
│   └── T001-task.md
└── T001-task-directa.md
```

## Prohibición de fallback

Nunca representar múltiples tasks como una lista dentro de un único archivo:

```text
<roadmap-root>/algo-tasks.md
```

Eso no es una materialización válida del roadmap.

Si una operación pretende crear N tasks, debe crear N archivos `TXXX-*.md`.
Si las tasks pertenecen a un Outcome, debe existir también el README del Outcome.

Si falta schema, `.stem`, `rootline`, permisos o estructura para crear archivos
canónicos, detenerse. No usar `Write` directo para inventar una estructura
alternativa.

## Materialización determinística

La ruta primaria para crear archivos del roadmap es:

```bash
roadmapctl materialize --plan <plan-json> --dry-run --repo <repo-path> --roadmap-root <roadmap-root> --output json
roadmapctl materialize --plan <plan-json> --apply --repo <repo-path> --roadmap-root <roadmap-root> --output json
```

El skill no debe duplicar numbering, `rootline new`, writes, actualización de tablas ni escritura de `blocked_by`; debe producir plan estructurado, revisar dry-run y delegar en `roadmapctl materialize`.

## Auto-numbering

El skill no calcula números `OXX`/`TXXX`. `roadmapctl materialize` asigna numbering determinístico y reporta las rutas propuestas en `changes[]` durante dry-run. Si el dry-run no produce rutas canónicas, detenerse y reportar diagnostics.

## Verificación de padre

Antes de crear un archivo, verificar que el directorio destino existe:

```bash
rootline describe <directorio>/
```

Si no existe, informar al usuario y no crear archivos fuera del roadmap.

Excepción permitida: `plan-subcommand.md` puede crear `<roadmap-root>/` y
`<roadmap-root>/.stem` durante su bootstrap explícito. Fuera de ese flujo, no
crear directorios ad-hoc.

## Cascading links

El skill no edita manualmente la tabla `## Tasks`. `roadmapctl materialize` actualiza el README del Outcome y mantiene la tabla sin columna Estado; el estado se lee desde frontmatter.

## Dependencias

Declarar `blocked_by` en la task bloqueada, con path relativo explícito.

- Misma carpeta/Outcome: `[[blocked_by:./T001-prerequisite.md]]`
- Otro Outcome: `[[blocked_by:../O01-setup/T001-prerequisite.md]]`
- No usar targets bare como `[[blocked_by:T001-prerequisite]]`; rootline solo los resuelve por basename único y pueden romperse con duplicados.

## Comandos Rootline de Referencia (troubleshooting/legacy)

| Comando | Cuándo usarlo |
|---------|---------------|
| `rootline validate <path>` | Después de crear/editar `.md` |
| `rootline fix <path>` | Si validate falla y la propuesta es segura |
| `rootline query <path> --where "expr"` | Listar tasks por metadata |
| `rootline tree <path> --where "expr" --output json` | Vista jerárquica con conteos |
| `rootline graph <path> --where "expr" --check` | Validar dependencias |

No usar `rootline stats`; `tree` ya incluye conteos.
