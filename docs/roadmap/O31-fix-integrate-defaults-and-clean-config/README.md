---
tipo: outcome
---
# O31: Fix `/integrate` defaults y limpieza de config muerta

Arreglar tres defectos relacionados en la cadena `roadmapctl` + `/integrate`:

1. **`pr_mode=false` (default del TOML) produce commits huérfanos**: el skill `/integrate` crea `feat/<scope>` siempre, también cuando `pr_mode=false`. Sin Fase 5/6 que mergee de vuelta, los commits quedan en el feature branch y master nunca los recibe. Fix: gate el `checkout -B` a `pr_mode=true`; con `pr_mode=false` el skill sincroniza con `base_branch` y commitea ahí.

2. **`pr_merge_strategy` es config muerta**: el campo está plumbed end-to-end (TOML → Config → JSON bootstrap → docs) pero ningún consumidor lo aplica (`gh pr merge --auto` deja decidir a GitHub). El warning de "deprecación gradual" no resuelve nada. Fix: eliminar completamente, sin migración ni fallback.

3. **`branch_style` fallback silencioso**: el skill inventa `feat/<scope>` si `branch_style` está vacío. Fix: emitir `RMC_INTEGRATE_BRANCH_STYLE_MISSING` y detener — el operador corre el wizard de gitflow.

Cuando este Outcome esté cerrado: el default `pr_mode=false` funciona sin sorpresas, la superficie de config sólo contiene campos con consumidor real, y la falta de `branch_style` falla temprano en vez de generar branches inventados.
