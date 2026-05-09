---
estado: Pending
tipo: task
---
# T005: Agregar fixtures/goldens read-only

**Outcome**: [O04 Estado read-only](README.md)
**Contribuye a**: CE1, CE2, CE3

[[blocked_by:./T002-implement-pending-command.md]]
[[blocked_by:./T003-implement-next-command.md]]
[[blocked_by:./T004-implement-decision-command.md]]

## Preserva

- INV1: Goldens son determinísticos y normalizan paths absolutos.
  - Verificar: `go test ./...`.

## Contexto

Los nuevos comandos read-only necesitan fixtures propios para direct tasks, outcomes, blockers, no pending y workspace.

## Alcance

**In**:
1. Fixtures para pending con direct tasks y outcomes.
2. Fixture con ready/blocked mixed.
3. Fixture no pending.
4. Fixture decision con reverse dependencies.
5. Goldens JSON/text relevantes.

**Out**:
- Fixtures de transition/materialize.

## Estado inicial esperado

- Comandos read-only implementados.

## Criterios de Aceptación

- Tests cubren `pending`, `next`, `decision` con fixtures.
- Golden output incluye `kind` correcto por comando.
- CI cross-platform pasa.

## Fuente de verdad

- `testdata/fixtures/*`
- `testdata/golden/*`
- `internal/cli/golden_test.go`
