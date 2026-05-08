---
estado: Completed
tipo: task
---
# T002: Crear esqueleto Go cross-platform para roadmapctl

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE1 y CE2

[[blocked_by:./T001-define-cli-contract.md]]

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: no se agrega código bajo `cmd/rootline`.
- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: solo existen comandos `doctor` y `check`.

## Contexto

El proyecto debe compilar en Linux, macOS y Windows. La estructura debe permitir crecer como suite de herramientas alrededor de Rootline sin mezclar lógica de roadmap en Rootline core.

## Alcance

**In**:
1. Inicializar módulo Go para `roadmapctl`.
2. Crear `cmd/roadmapctl/main.go`.
3. Crear paquetes internos mínimos: `internal/cli`, `internal/diagnostics`, `internal/config`, `internal/rootlinecli`, `internal/roadmap`, `internal/fsx`.
4. Configurar Cobra o parser equivalente para `doctor` y `check`.
5. Añadir `go test ./...` y `go build ./cmd/roadmapctl` como comandos esperados.

**Out**:
- Implementar lógica completa de checks.
- Configurar release completo.
- Añadir comandos fuera de `doctor`/`check`.

## Estado inicial esperado

- Contrato CLI aprobado o documentado.
- Repo contiene roadmap y README inicial.

## Criterios de Aceptación

- `go test ./...` pasa.
- `go build ./cmd/roadmapctl` genera binario local.
- `roadmapctl --help`, `roadmapctl doctor --help` y `roadmapctl check --help` funcionan.
- No hay imports de `github.com/pablontiv/rootline/internal/*`.

## Fuente de verdad

- `go.mod`
- `cmd/roadmapctl/main.go`
- `internal/cli/`
- `README.md`
