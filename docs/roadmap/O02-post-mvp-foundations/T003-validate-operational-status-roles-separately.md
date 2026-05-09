---
estado: Pending
tipo: task
---
# T003: Validar roles operacionales separados del schema

**Outcome**: [O02 Fundaciones post-MVP](README.md)
**Contribuye a**: CE2, INV1

[[blocked_by:./T002-make-stem-authoritative-for-document-schema.md]]

## Preserva

- INV1: Un valor documental válido por `.stem` no se invalida por no ser rol operacional.
  - Verificar: fixture `valid-status-on-hold`.

## Contexto

`status-values`, `done-statuses` y `active-statuses` son configuración operacional. Deben mapear conceptos como completed/in-progress/done, pero no reducen el conjunto de valores documentales permitidos.

## Alcance

**In**:
1. Crear validación separada para roles config contra `schema.estado.values`.
2. Emitir diagnostic específico si un role apunta a un estado inexistente en schema.
3. Validar `done-statuses` y `active-statuses` como listas operacionales.
4. Documentar que config roles no son enum exhaustivo.

**Out**:
- Cambiar nombres concretos de estados por defecto.
- Implementar transiciones.

## Estado inicial esperado

- T002 ya separó validación documental de schema.

## Criterios de Aceptación

- Config con `completed = Done` y schema sin `Done` falla con diagnostic de config/schema mismatch.
- Task con `estado: On Hold` y schema válido pasa aunque `On Hold` no sea role.
- Tests cubren role inválido, active/done inválidos y config válida.

## Fuente de verdad

- `internal/config/config.go`
- `internal/roadmap/status.go`
- `docs/cli-contract.md`
- `testdata/fixtures/invalid-config-role-not-in-schema`
