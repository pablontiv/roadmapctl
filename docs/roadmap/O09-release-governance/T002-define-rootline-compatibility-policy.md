---
estado: Completed
tipo: task
---
# T002: Definir política de compatibilidad Rootline

**Outcome**: [O09 Release/governance](README.md)
**Contribuye a**: CE2, INV1

[[blocked_by:../O02-post-mvp-foundations/T004-harden-rootline-json-on-nonzero-exit.md]]

## Preserva

- INV1: No hacer hard-fail por versión sin decisión explícita.
  - Verificar: docs/diagnostics.

## Contexto

La sesión decidió no fijar todavía hard-fail de versión. Aun así, los comandos post-MVP dependerán de `tree`, `set`, `new`, `describe` shapes y JSON estable.

## Alcance

**In**:
1. Detectar y reportar versión Rootline.
2. Definir minimum/recommended version como warning o error según decisión.
3. Agregar tests/fakes para comandos requeridos.
4. Actualizar CI para latest y/o versión mínima si se aprueba.

**Out**:
- Cambiar Rootline.
- Bloquear usuarios por versión sin política aprobada.

## Estado inicial esperado

- `doctor` reporta `rootline --version` como info.

## Criterios de Aceptación

- Docs explican compatibilidad requerida por comando.
- Diagnostics diferencian missing binary, incompatible command y JSON shape inválido.
- CI cubre al menos versión instalada actual.

## Fuente de verdad

- `docs/release.md`
- `internal/rootlinecli/*`
- `.github/workflows/ci.yml`
