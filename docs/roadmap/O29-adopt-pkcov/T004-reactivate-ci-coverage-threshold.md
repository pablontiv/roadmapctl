---
estado: Specified
tipo: task
---
# T004: Reactivar coverage-threshold en CI (0 → 85)

**Outcome**: [O29 Adoptar pkcov + reactivar gate](README.md)
**Contribuye a**: roadmapctl entra en el patrón compartido con rootline/backscroll (INV1)

[[blocked_by:./T001-raise-cmd-roadmapctl-coverage.md]]
[[blocked_by:./T002-cover-or-remove-internal-templates.md]]
[[blocked_by:./T003-raise-internal-workspace-coverage.md]]

## Preserva

- INV1 del outcome: threshold uniforme 85
- INV3 del outcome: el smoke job sigue corriendo en paralelo — validan cosas distintas
  - Verificar: smoke job sigue habilitado en `ci.yml` post-cambio

## Contexto

`/home/shared/roadmapctl/.github/workflows/ci.yml:18` declara `coverage-threshold: 0  # cobertura validada por smoke job en 3 plataformas`. El comentario sugiere que el smoke job (líneas 51-53, corriendo `./scripts/check-coverage.sh` en Ubuntu/macOS/Windows) sustituye al gate de coverage del workflow base de crossbeam. Esto no es preciso: el smoke job valida que el binario corra en 3 OS, no la cobertura del código.

Con T001-T003 completos, todos los paquetes pasan ≥85 y el total también. Reactivar el threshold ahora no rompe nada.

## Alcance

**In**:

1. Editar `/home/shared/roadmapctl/.github/workflows/ci.yml`: cambiar `coverage-threshold: 0` a `coverage-threshold: 85`. Actualizar o eliminar el comentario inline (era engañoso).
2. Mantener el smoke job intacto — sigue validando cross-platform build, ortogonal al gate de coverage.
3. Verificar que CI verde tras push: el gate de crossbeam debe pasar con cobertura ≥85.

**Out**:
- No tocar el smoke job.
- No cambiar `scripts/check-coverage.sh` (T005 lo reemplaza).
- No bumpear el threshold por encima de 85 (este task es restaurar paridad).

## Estado inicial esperado

- T001, T002, T003 completadas: todos los paquetes ≥85, total ≥85.
- `ci.yml:18` declara `coverage-threshold: 0`.

## Criterios de Aceptación

- `.github/workflows/ci.yml` declara `coverage-threshold: 85`.
- Comentario inline removido o actualizado para reflejar la verdad.
- CI verde post-push (gate de coverage pasa).
- Smoke job sigue habilitado y verde.

## Fuente de verdad

- `/home/shared/roadmapctl/.github/workflows/ci.yml`
