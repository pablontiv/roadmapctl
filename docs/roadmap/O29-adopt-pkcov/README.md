---
tipo: outcome
---
# O29: Adoptar pkcov de picokit + reactivar gate de coverage

Roadmapctl hoy tiene el gate de coverage **desactivado** en CI (`coverage-threshold: 0` en `ci.yml:18`, deferred al smoke job). Esto es una excepción al patrón compartido con rootline/backscroll. Resultado observable: tres paquetes muy débiles — `cmd/roadmapctl` 50%, `internal/templates` 0% (GenerateStemContent sin tests), `internal/workspace` 82.6% — pasan inadvertidos. El total reporta 86.2% pero esos paquetes están bajo el floor común.

Este outcome trae roadmapctl al patrón compartido: pre-work para subir los 3 paquetes débiles a ≥85, reactivar el gate de CI a 85 (igualando rootline/backscroll), y adoptar `pkcov` de picokit (que ya está en `go.mod` por `diag`/`pathsec`/`autoupdate`) reemplazando el `scripts/check-coverage.sh` actual (30+ líneas bash + python3).

Es el caso más urgente de los tres consumidores: el gate desactivado significa que cualquier regresión actual o futura en roadmapctl no se detecta hasta que alguien la note manualmente. El pre-work además limpia paquetes que llevan tiempo sin atención.

Invariantes preservadas:
- INV1: threshold uniforme 85 (mismo que rootline/backscroll/picokit)
- INV2: el pre-work no relaja el contrato — sube los paquetes débiles, no agrega excepciones
- INV3: la decisión de mantener el smoke job separado se documenta — coverage gate y smoke job validan cosas distintas

Scope: pre-work moderado + reactivación + adopción mecánica. No introduce features nuevas; no cambia la lógica de roadmapctl. Depende cross-repo de picokit `O03-coverage-tooling` Completed con tag publicado para el swap de bash → pkcov, pero el pre-work + reactivación pueden empezar en paralelo a O03.
