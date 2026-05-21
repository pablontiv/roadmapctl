---
estado: Specified
tipo: task
---
# T006: Agregar pre-push gate + docs

**Outcome**: [O29 Adoptar pkcov + reactivar gate](README.md)
**Contribuye a**: roadmapctl detecta regresiones antes de CI; declara conformance con coverage-spec v1.0

[[blocked_by:./T005-add-floors-and-replace-bash.md]]

## Preserva

- INV1 del outcome: threshold 85 aplicado antes del push
  - Verificar: regresión simulada bloquea `git push`

## Contexto

coverage-spec v1.0 sección 5 requiere que el repo ejecute `pkcov check` desde `.githooks/pre-push` cuando cambian archivos `*.go`. El patrón a imitar vive en `/home/shared/rootline/.githooks/pre-push`.

## Alcance

**In**:

1. Editar `/home/shared/roadmapctl/.githooks/pre-push` (o crear si no existe) agregando bloque condicional sobre `*.go` que llama `just coverage-check`.
2. Verificar que git hooks están enlazados (`git config core.hooksPath`); si no, documentar en CLAUDE.md o README.
3. Actualizar CLAUDE.md o README:
   - Sección de comandos: `just coverage` / `just coverage-check`
   - Sección de CI: referencia a `picokit/docs/coverage-spec.md`, declarar conformance v1.0
4. Test: simular regresión (`git rm <test>.go && git commit && git push`) → debe bloquearse; restaurar.

**Out**:
- No cambiar comportamiento de CI workflow (lo activó T004).
- No documentar internals de pkcov.

## Estado inicial esperado

- T005 completada: `just coverage-check` invoca pkcov y exit 0.
- `.githooks/pre-push` no menciona coverage.

## Criterios de Aceptación

- `.githooks/pre-push` incluye bloque de coverage gate condicional sobre `*.go`.
- Simulación de regresión bloquea push.
- Push docs-only (`*.md`) no dispara el gate.
- README o CLAUDE.md tiene sección de coverage que referencia el spec.
- Declaración "roadmapctl cumple coverage-spec v1.0" en el doc.

## Fuente de verdad

- `/home/shared/roadmapctl/.githooks/pre-push`
- `/home/shared/roadmapctl/README.md` o `CLAUDE.md`
- `/home/shared/rootline/.githooks/pre-push` — patrón
- `/home/shared/picokit/docs/coverage-spec.md` — spec
