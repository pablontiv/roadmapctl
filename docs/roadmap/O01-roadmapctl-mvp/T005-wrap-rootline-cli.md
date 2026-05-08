---
estado: Completed
tipo: task
---
# T005: Encapsular llamadas seguras al binario rootline

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE1 y CE2

[[blocked_by:./T002-create-go-cli-skeleton.md]]

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: se usa Rootline como proceso externo, no imports de internals.
- INV2: `roadmapctl` no invoca subprocess con shell strings.
  - Verificar: solo `exec.CommandContext` con args explícitos.

## Contexto

`roadmapctl` necesita usar Rootline para validar `.stem`, consultar registros y analizar grafos, pero debe tratarlo como dependencia externa con contrato JSON y control de errores.

## Alcance

**In**:
1. Crear `internal/rootlinecli.Client`.
2. Resolver binario por `--rootline`, `ROOTLINE_BIN` o PATH.
3. Ejecutar con `exec.CommandContext`, timeout y args explícitos.
4. Controlar `Dir` y `Env` mínimos.
5. Capturar stdout/stderr separadamente.
6. Parsear JSON para `validate`, `describe`, `query` y `graph`.
7. Convertir errores de entorno a exit `3`.

**Out**:
- Usar shell, pipes o redirecciones.
- Importar `github.com/pablontiv/rootline/internal/*`.
- Implementar checks de negocio.

## Estado inicial esperado

- Skeleton CLI existe.
- Contrato de diagnostics existe o está en progreso.

## Criterios de Aceptación

- Tests unitarios usan fake executor o fake client.
- Missing `rootline` produce diagnostic claro.
- Timeout produce diagnostic claro y exit `3`.
- Comandos generados no pasan por shell.
- JSON inválido de Rootline produce error controlado.

## Fuente de verdad

- `internal/rootlinecli/`
- `internal/diagnostics/`
- `docs/cli-contract.md`
