---
tipo: outcome
---
# O27: Complete picokit integration (autoupdate + facade cleanup)

Cuando todas las tasks estén completadas, roadmapctl habrá completado la adopción de picokit eliminando `internal/updater/` (paquete origen del actual `picokit/autoupdate`) y la facade de `internal/diagnostics/` que hoy solo re-exporta tipos vía type alias.

Resultado esperado: codebase ~1554 líneas menor, código de auto-update vive en una sola fuente upstream (`github.com/pablontiv/picokit/autoupdate`), y los tipos de diagnóstico se importan directamente de `picokit/diag` con un paquete mínimo `internal/reports` que conserva el campo `RoadmapRoot` específico de roadmapctl.

Continúa el trabajo de O25 (que migró `pathsec`, `diag` vía facade y eliminó `internal/diff`). El motivador original de picokit era `autoupdate`; O25 lo dejó pendiente. O27 lo cierra.

**Prerequisito externo:** picokit v0.1.1 (o superior) debe estar disponible vía `go get`. Hoy ya está taggeado y resolvible.
