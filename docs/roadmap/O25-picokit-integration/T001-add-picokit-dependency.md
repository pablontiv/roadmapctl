---
estado: Completed
tipo: task
---
# T001: Add picokit@v0.1.0 dependency to roadmapctl

**Outcome**: [O25 Integrate picokit as a dependency](README.md)
**Contribuye a**: picokit está disponible como dependencia importable en el módulo.

## Preserva

- INV1: `go build ./...` pasa sin errores tras agregar la dependencia.
  - Verificar: `go build ./...` desde `/home/shared/roadmapctl`.
- INV2: `go test ./...` pasa sin regresiones.
  - Verificar: `go test ./...` desde `/home/shared/roadmapctl`.

## Contexto

picokit v0.1.0 es el módulo `github.com/pablontiv/picokit`. Esta task solo agrega la dependencia al módulo; las tasks siguientes (T002-T004) realizan las migraciones reales usando los paquetes de picokit.

**Prerequisito externo**: El tag `v0.1.0` de picokit debe estar pusheado al remote. Si no está disponible, esta task está bloqueada externamente.

## Alcance

**In**:
1. `go get github.com/pablontiv/picokit@v0.1.0` desde `/home/shared/roadmapctl`.
2. `go mod tidy` para limpiar go.sum.

**Out**:
- No modificar ningún archivo `.go` de código fuente en esta task (solo go.mod y go.sum).
- No eliminar ningún paquete interno aún (eso lo hacen T002-T004).

## Estado inicial esperado

- `grep picokit /home/shared/roadmapctl/go.mod` retorna vacío.
- picokit v0.1.0 está disponible en el registry de GitHub (`GOPROXY=direct go get github.com/pablontiv/picokit@v0.1.0` resuelve).

## Criterios de Aceptación

- `grep "pablontiv/picokit" /home/shared/roadmapctl/go.mod` muestra `v0.1.0`.
- `go build ./...` pasa desde `/home/shared/roadmapctl`.
- `go test ./...` pasa desde `/home/shared/roadmapctl`.

## Fuente de verdad

- `/home/shared/roadmapctl/go.mod`
- `/home/shared/roadmapctl/go.sum`
