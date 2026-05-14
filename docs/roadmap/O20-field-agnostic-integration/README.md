---
tipo: outcome
---
# O20: Field-agnostic integration

roadmapctl hardcodea nombres de campo (`"estado"`, `"tipo"`, `"blocked_by"`) cuando parsea resultados de rootline. Esto crea acoplamiento entre el controller y el schema del roadmap.

Este outcome mueve ese vocabulario a `.roadmapctl.toml` bajo una sección `[fields]` con defaults retrocompatibles. Todos los consumers internos (model, status, next, dependencies, structure, schema_portability, bootstrap) dejan de hardcodear nombres de campo y usan `cfg.Fields.*` en su lugar.

El resultado observable es que roadmapctl funciona correctamente con cualquier nombre de campo configurado en `[fields]`, y el comportamiento por defecto es idéntico al actual (defaults retrocompatibles garantizan que repos existentes no necesitan cambiar su config).

Coordina con el outcome [O14 en rootline](../../rootline/O14-field-agnostic-refactor/README.md) que elimina los hardcodings del engine.
