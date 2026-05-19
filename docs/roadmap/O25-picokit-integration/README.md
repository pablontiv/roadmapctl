---
tipo: outcome
---
# O25: Integrate picokit as a dependency

roadmapctl adopta picokit como dependencia explícita, eliminando código utilitario duplicado que fue la fuente original de los paquetes de picokit.

El resultado es un codebase más pequeño (~400 líneas eliminadas) donde los utilitarios genéricos viven en un único lugar canónico. Los paquetes migrados son: `pathsec` (reemplaza `internal/fsx`), `diag` (reemplaza la infraestructura de `internal/diagnostics`). El paquete `internal/diff` es dead code que se elimina directamente.

**Prerequisito externo**: picokit v0.1.0 debe estar tageado y resolvible vía `go get` antes de ejecutar T001 de este outcome.
