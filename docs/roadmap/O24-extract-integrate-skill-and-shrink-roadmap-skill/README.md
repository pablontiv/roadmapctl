---
tipo: outcome
---
# Extraer skill `/integrate` y achicar el skill `/roadmap`

El skill `/roadmap` mezcla tres responsabilidades — planificación, ejecución del
loop y gitflow (commit / push / branch / PR / merge). El resultado son ~1458
líneas en 11 archivos, con prosa de gitflow en `loop-subcommand.md` (~30 líneas
del paso 9) y en `pr-workflow.md` (97 líneas), cargada cada vez que el loop
corre. Además, `SKILL.md` carga 312 líneas en cualquier invocación de `/roadmap`,
incluso para subcomandos como `pending` que solo necesitan ~40.

Al terminar este Outcome existirá un skill aparte `/integrate` que encapsula
toda la prosa de gitflow per-task — branch por scope, commit con `commit_style`,
push si `auto_push`, `gh pr create/merge` según `pr_mode`/`autonomy`/`pr_merge_strategy`,
cleanup post-merge — y será invocable tanto desde el paso 9 de `/roadmap loop`
como ad-hoc. En paralelo, `SKILL.md` del skill `/roadmap` quedará ≤ 100 líneas
gracias a progressive disclosure: bootstrap detallado, tabla completa de config
y notas de contributor mudarán a `*-reference.md` cargables on-demand.

Diff esperado: skill `/roadmap` pasa de ~1458 a ~900 líneas (`pr-workflow.md`
eliminado, `SKILL.md` 312 → ≤ 100); skill `/integrate` nuevo aporta ~120 líneas
reusables. No se reimplementa `git`/`gh` en Go — la prescripción de comandos
sigue siendo prosa, pero modular.

Plan completo en `/home/pones/.claude/plans/tiene-sentido-que-roadmap-modular-stroustrup.md`.
