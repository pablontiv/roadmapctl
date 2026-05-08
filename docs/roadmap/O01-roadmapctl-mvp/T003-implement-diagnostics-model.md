---
estado: Completed
tipo: task
---
# T003: Implementar modelo de diagnostics y renderers

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE1 y CE2

[[blocked_by:./T001-define-cli-contract.md]]

## Preserva

- INV2: `roadmapctl` no invoca subprocess con shell strings.
  - Verificar: esta task no agrega ejecución de comandos externos.
- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: diagnostics solo reportan problemas.

## Contexto

`roadmapctl check` y `roadmapctl doctor` deben ser útiles para humanos y automatización. El formato JSON estable es la interfaz que el skill `/roadmap`, CI y hooks podrán consumir.

## Alcance

**In**:
1. Definir tipos `Report`, `Summary`, `Diagnostic`, `Severity`.
2. Implementar renderer JSON a stdout sin logs mezclados.
3. Implementar renderer text/human-readable.
4. Definir helpers para contar errores/warnings y calcular exit code.
5. Añadir diagnostic IDs iniciales documentados.

**Out**:
- Checks de roadmap.
- Invocación a Rootline.
- Integración con CI.

## Estado inicial esperado

- Existe skeleton Go.
- Existe contrato JSON/exit-code.

## Criterios de Aceptación

- Tests unitarios cubren salida JSON para report success/fail.
- `--output json` produce JSON parseable sin texto adicional en stdout.
- Errores/logs van a stderr o se suprimen en modo JSON.
- Exit code derivado coincide con el contrato: `0`, `1`, `2`, `3`, `4`.

## Fuente de verdad

- `internal/diagnostics/`
- `internal/cli/`
- `docs/cli-contract.md`
