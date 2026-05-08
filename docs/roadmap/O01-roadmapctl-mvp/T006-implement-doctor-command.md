---
estado: Completed
tipo: task
---
# T006: Implementar roadmapctl doctor

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE1

[[blocked_by:./T003-implement-diagnostics-model.md]]
[[blocked_by:./T004-load-roadmap-config.md]]
[[blocked_by:./T005-wrap-rootline-cli.md]]

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: doctor diagnostica Rootline, no modifica Rootline.
- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: doctor no escribe archivos.

## Contexto

`doctor` debe explicar por qué un comando `/roadmap` puede o no puede operar. Es el preflight obligatorio para comandos implementados que escriben, mutan o ejecutan.

## Alcance

**In**:
1. Detectar repo o workspace desde `--repo` o cwd.
2. Detectar `rootline` y versión disponible.
3. Leer configuración de roadmap.
4. Verificar presencia de `<roadmap-root>` y `.stem`.
5. Reportar rutas config/cache relevantes si existen.
6. Emitir JSON y text con diagnostics accionables.

**Out**:
- Ejecutar checks profundos de estructura.
- Modificar hooks.
- Instalar dependencias automáticamente.

## Estado inicial esperado

- Config loader y rootline client existen.
- Diagnostics renderer existe.

## Criterios de Aceptación

- `roadmapctl doctor --repo testdata/fixtures/valid-outcome-with-tasks --output json` devuelve status ok.
- Fixture sin Rootline configurable devuelve diagnostic de entorno y exit `3`.
- Fixture sin `.claude/roadmap.local.md` reporta config faltante.
- Salida JSON es parseable y no mezcla logs en stdout.

## Fuente de verdad

- `internal/cli/doctor.go`
- `internal/config/`
- `internal/rootlinecli/`
- `testdata/fixtures/`
