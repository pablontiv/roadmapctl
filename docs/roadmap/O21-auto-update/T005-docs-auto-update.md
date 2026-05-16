---
estado: Completed
tipo: task
---
# T005: Documentar el feature y patrón de auto-update

**Outcome**: [O21 Auto-update staged async](README.md)
**Contribuye a**: Que futuros contribuidores entiendan el patrón staged async sin leer el código

[[blocked_by:./T003-cli-wiring.md]]

## Preserva

- INV1: `docs/cli-contract.md` no documenta comportamiento de actualización que no esté implementado.
  - Verificar: Cada AC documentado tiene cobertura en los tests de T004.

## Contexto

El patrón staged async es no obvio: el update se descarga en background en el run N y se aplica en el run N+1 con re-exec. Sin documentación, esto puede parecer un bug (la versión no cambia inmediatamente) o generar confusión sobre cuándo aplica y cuándo no.

Documentar en dos lugares:
1. `docs/` — documento técnico del patrón (cómo funciona, escape hatches, comportamiento en cada OS)
2. `README.md` del repo — sección usuario-facing (cómo saber si hay un update, cómo forzar/saltear)

## Alcance

**In**:
1. Crear `docs/auto-update.md` con: descripción del patrón staged async, flujo por invocación, comportamiento por OS (Unix/Windows), escape hatches (`ROADMAPCTL_NO_UPDATE=1`, `version == "dev"`), qué hacer si el update falla silenciosamente
2. Agregar sección "Auto-update" en `README.md` (usuario-facing, sin detalles internos)

**Out**:
- No modificar `docs/cli-contract.md` (el auto-update no es parte del contrato de CLI estable)
- No documentar internals de Go ni el paquete `internal/updater` (eso es responsabilidad de los comentarios en código)

## Estado inicial esperado

- T003 completada: el feature está integrado y funcional

## Criterios de Aceptación

- `docs/auto-update.md` existe y cubre: flujo por invocación, lag de 1 run, comportamiento Unix vs Windows, escape hatches, qué pasa si hay error de permisos
- `README.md` tiene sección "Auto-update" con al menos: cómo funciona en una oración, cómo desactivarlo con `ROADMAPCTL_NO_UPDATE=1`
- Ningún AC documenta comportamiento que no esté cubierto por los tests de T004

## Fuente de verdad

- `docs/auto-update.md` (crear)
- `README.md` (modificar — agregar sección)
- `internal/updater/updater.go` y `apply.go` (fuente de verdad del comportamiento real)
