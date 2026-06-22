# Research: Provider-neutral integration context

**Fecha:** 2026-05-22
**Status:** Research
**Trigger:** evaluar mover logica de skills a `roadmapctl` sin acoplar gitflow a un proveedor.

---

## Problema

Los skills contienen logica repetible de integracion post-task: detectar si una
task completada pertenece a un Outcome, decidir si es la ultima del grupo,
pasar configuracion de commit/push/PR, y evitar que el loop recalcule estado
del roadmap en prompt.

Mover esa logica completa a `roadmapctl` seria riesgoso si implica ejecutar
Git, crear PRs o asumir GitHub. La herramienta se usa en ambientes con GitHub,
Azure DevOps, GitLab y flujos locales. El core debe permanecer portable.

El problema no es que `roadmapctl` sepa hechos del roadmap. El problema seria
que convierta esos hechos en politica gitflow o en acciones contra un proveedor.

## Frontera acordada

`roadmapctl` entrega contexto factual y validado sobre la task completada y la
config efectiva. No ejecuta comandos de integracion.

`/integrate` conserva la responsabilidad de adaptar ese contexto al ambiente:
GitHub, Azure DevOps, GitLab, local/direct-push u otro proveedor. Tambien
conserva la generacion de nombres y texto cuando dependen de convenciones del
repo.

La nomenclatura del roadmap (`OXX`, `TXXX`, slugs de Outcome) es contexto. No es
politica gitflow. Solo debe influir en branch, commit o PR si el estilo del repo
lo pide explicitamente.

## No-goals

- No agregar `roadmapctl integrate plan`.
- No ejecutar `git`, `gh`, `az repos`, `glab` ni CLIs de proveedor desde
  `roadmapctl`.
- No crear, cerrar ni mergear PRs desde `roadmapctl`.
- No devolver `branch_name`, `branch_target` ni `branch_template` calculados
  desde la nomenclatura del roadmap.
- No convertir automaticamente `OXX`, `TXXX`, `outcome_slug` o `task_id` en
  nombres de branch.
- No reintroducir `roadmapctl materialize` ni un writer semantico de Markdown.

## Contrato propuesto

Extender `roadmapctl transition can-complete` y
`roadmapctl transition complete --apply` para incluir un bloque
`integration_context` provider-neutral.

Ejemplo:

```json
{
  "integration_context": {
    "task_path": "docs/roadmap/O31-x/T003-y.md",
    "task_id": "T003",
    "task_title": "...",
    "is_direct_task": false,
    "outcome_path": "docs/roadmap/O31-x/README.md",
    "outcome_id": "O31",
    "outcome_slug": "x",
    "outcome_title": "...",
    "is_last_task_in_outcome": false,
    "base_branch": "master",
    "auto_push": true,
    "pr_mode": false,
    "commit_style": "conventional",
    "branch_style": "",
    "pr_title_style": "",
    "pr_body_style": ""
  }
}
```

El bloque es factual:

- Identifica la task completada.
- Identifica si la task es directa o pertenece a un Outcome.
- Indica si es la ultima task activa de ese Outcome.
- Expone la configuracion efectiva necesaria para que el adaptador integre.
- No prescribe comandos ni nombres de branch.

## Semantica de `branch_style` vacio

`branch_style = ""` significa integrar directo a `base_branch`.

Consecuencias:

- No dispara wizard de gitflow.
- No genera branch.
- No inventa `feat/<scope>`, `feat/<outcome>` ni variantes derivadas del roadmap.
- `/integrate` debe hacer checkout/pull de `base_branch`, commitear ahi, y pushear
  `base_branch` si `auto_push = true`.

`pr_mode = true` con `branch_style = ""` es una configuracion contradictoria:
PR mode requiere una rama o convencion explicita desde donde abrir PR. En ese
caso el adaptador debe detenerse con un diagnostic claro y pedir configurar
`branch_style`.

`base_branch` debe ser configuracion explicita, no heuristica. Si falta,
`roadmapctl` debe reportar diagnostic de config antes de producir un contexto de
integracion util.

## Responsabilidades

### `roadmapctl`

- Validar roadmap y config.
- Mutar estado solo via `transition`.
- Calcular hechos del roadmap: task directa vs Outcome, metadata de Outcome,
  ultima task activa del Outcome.
- Exponer `integration_context` en JSON estable.
- Permanecer provider-neutral.

### `/roadmap loop`

- Ejecutar bootstrap, preflight, ACs y transiciones.
- Consumir `integration_context` desde `transition complete`.
- Pasar ese bloque a `/integrate`.
- No recalcular `is_last_task_in_outcome`.
- No derivar politica gitflow desde paths del roadmap.

### `/integrate`

- Interpretar `branch_style`, `commit_style`, `pr_title_style` y `pr_body_style`.
- Generar branch, commit message, PR title/body solo segun convenciones del repo.
- Ejecutar el proveedor disponible: GitHub, Azure DevOps, GitLab, local/direct.
- Usar datos del roadmap solo como contexto, no como regla automatica de naming.

### `/retrospective`

- Permanece como analisis de comprension y ejecucion.
- Puede proponer cambios a skills o docs, pero no necesita logica nueva en
  `roadmapctl`.

## Riesgos

- **Acoplar roadmap naming al gitflow:** introduciria ramas con `OXX/TXXX` aunque
  el repo no use esa nomenclatura.
- **Provider lock-in:** ejecutar `gh` desde `roadmapctl` excluiria ambientes de
  Azure DevOps, GitLab o flujos locales.
- **Duplicacion de politica:** si el skill y `roadmapctl` calculan ramas o ultima
  task por separado, los incidentes vuelven a aparecer.
- **Config ambigua:** `base_branch` no debe depender de deteccion implicita si se
  usa para integracion automatizada.

## Proximo paso

Disenar tasks de implementacion separadas:

1. Agregar `base_branch` explicito a `[gitflow]`, bootstrap y contrato CLI.
2. Agregar `integration_context` a `transition can-complete` y
   `transition complete --apply`.
3. Actualizar `/roadmap loop` para pasar `integration_context` a `/integrate`.
4. Actualizar `/integrate` para tratar `branch_style = ""` como direct-push a
   `base_branch` y no inventar ramas.
5. Agregar pruebas Go y headless de skills para cerrar la frontera.
