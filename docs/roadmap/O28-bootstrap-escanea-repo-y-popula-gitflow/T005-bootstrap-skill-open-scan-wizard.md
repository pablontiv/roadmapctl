---
estado: Specified
tipo: task
---
# T005: bootstrap-reference.md + nuevo bootstrap-subcommand.md

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: el skill describe escaneo abierto del repo y wizard pre-llenado

[[blocked_by:./T002-bootstrap-serialize-and-diagnostic.md]]

## Preserva

- INV1: bootstrap-reference.md sigue siendo referenciado por los subcomandos existentes
  - Verificar: `grep -r 'bootstrap-reference' .claude/skills/roadmap/` retorna resultados

## Contexto

Cuando bootstrap devuelve `RMC_GITFLOW_NOT_CONFIGURED`, el skill `/roadmap` debe:
1. Escanear el repo completo (no una lista fija de comandos) para entender las convenciones
2. Redactar propuesta de TOML con los 3 style fields + campos determinísticos
3. Presentar wizard pre-llenado campo por campo (usuario confirma/edita/salta)
4. Escribir el TOML solo tras confirmación final
5. Re-ejecutar bootstrap (consume TOML recién escrito) y continuar

El subcomando `/roadmap bootstrap` fuerza el re-escaneo on-demand aunque el TOML exista.

El escaneo es **abierto**: el LLM decide qué leer — todos los .md, CLAUDE.md del workspace y del proyecto, git log, gh pr list (merged + closed), branch protection, workflows, código fuente, cualquier otro archivo relevante. No hay lista fija.

## Alcance

**In**:
1. Actualizar `bootstrap-reference.md`: agregar sección que describe detección de `RMC_GITFLOW_NOT_CONFIGURED` → escaneo abierto → wizard → escritura TOML → re-bootstrap
2. Crear `.claude/skills/roadmap/bootstrap-subcommand.md`: describe `/roadmap bootstrap` on-demand (re-escaneo + wizard, termina sin ejecutar plan/loop/etc.)

**Out**:
- No modificar SKILL.md routing (T006)
- No modificar los otros subcomandos (T007)

## Estado inicial esperado

- `.claude/skills/roadmap/bootstrap-reference.md` existe
- T002 completada

## Criterios de Aceptación

- `bootstrap-reference.md` contiene sección sobre `RMC_GITFLOW_NOT_CONFIGURED` y describe escaneo abierto (sin lista fija de comandos)
- El texto de escaneo menciona explícitamente que el LLM decide qué leer según lo que encuentra
- `.claude/skills/roadmap/bootstrap-subcommand.md` existe y describe el subcomando on-demand
- Ninguno de los dos archivos contiene una lista rígida de exactamente N comandos a ejecutar

## Fuente de verdad

- `.claude/skills/roadmap/bootstrap-reference.md`
- `.claude/skills/roadmap/bootstrap-subcommand.md` (nuevo)
