---
estado: Completed
tipo: task
---

# T015 Fix transition start/complete: requerir --apply explícito

`roadmapctl transition start/complete <path>` sin `--apply` retorna exit 0 y `allowed: true` pero no aplica el cambio. Los agentes creen que la transición fue aplicada cuando no lo fue.

## Criterios de Aceptación

- AC1: `roadmapctl transition complete <path>` sin `--apply` → exit 2, diagnostic `RMC_TRANSITION_APPLY_FAILED`
- AC2: `roadmapctl transition start <path>` sin `--apply` → mismo comportamiento
- AC3: `roadmapctl transition can-start` y `can-complete` no se ven afectados (siguen siendo dry-run válidos)
- AC4: `go test ./...` pasa; test nuevo cubre el caso `start`/`complete` sin `--apply`

## Especificación Técnica

Archivo: `internal/cli/transition.go`

La condición actual de error es:
```go
if !dryRun && !apply { ... }  // solo dispara con --dry-run=false explícito
```

Cambiar para que `start` y `complete` requieran siempre `--apply`:
```go
if !apply && (action == "start" || action == "complete") {
    // retornar error RMC_TRANSITION_APPLY_FAILED
}
```

`can-start` y `can-complete` no entran en esta condición — siguen siendo dry-run por diseño.

Agregar test en `internal/cli/transition_test.go` que verifique exit 2 + `RMC_TRANSITION_APPLY_FAILED` para `transition complete <path>` sin `--apply`.

## Contexto

Causa raíz del bug en cartyx: 4 tasks quedaron en `In Progress` porque agentes corrieron `roadmapctl transition complete` sin `--apply`. El comando reportó éxito (exit 0, `allowed: true`) pero no escribió el status.
