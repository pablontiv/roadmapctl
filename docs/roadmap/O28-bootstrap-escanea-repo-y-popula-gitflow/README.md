---
tipo: outcome
---
# O28: Bootstrap escanea repo y popula gitflow

Cuando todas las tasks estén completadas, `/roadmap` inspeccionará el contexto real de cada repo antes de operar: escanea el repo completo (código, documentación, git, GitHub), popula `branch_style`, `pr_title_style` y `pr_body_style` en `.roadmapctl.toml`, y el LLM en `/integrate` genera branch names, commits y PRs siguiendo esas convenciones — sin hardcoding.

Corrige dos bugs documentados: (1) el loop podía saltarse `/integrate` y hacer git directo; (2) `pr_mode=false` instruía commitear en el branch actual aunque fuera `main`.
