---
estado: Specified
tipo: task
---
# T029: Loop — correr linter localmente antes de declarar ACs pasados en tasks con CI

**Contribuye a**: Failures de lint pre-existentes o introducidos por la task se detectan localmente antes del push, no después de que CI falla.

## Contexto

El linter en CI (golangci-lint, eslint, clippy, etc.) corre sobre el estado completo del repo, incluyendo archivos no tocados por la task activa. Si hay violations pre-existentes, el primer push con código fuente las expone — y la task que las descubre no es la responsable de introducirlas.

Ejemplo real: picokit tenía 6 violations de lint (errcheck, staticcheck SA5011/SA9003, unused) en paquetes no tocados por T002. T003 heredó esa deuda y tuvo que incluir 3 commits de fix no planificados.

## Alcance

**In**:
1. En `loop-subcommand.md`, Fase 3 paso 5 "Implementar" o paso 6 "Verificar ACs", agregar instrucción:
   - Si la task modifica código fuente (`.go`, `.ts`, `.js`, `.rs`, `.py`, etc.) y el repo tiene un linter configurado (`.golangci.yml`, `.eslintrc`, `clippy.toml`, etc.), **correr el linter sobre los paquetes/módulos afectados** antes de declarar ACs pasados.
   - Si el linter reporta violations en archivos **no tocados por la task activa**, documentarlas en el contexto de sesión como candidatas a nueva task de fix — no implementarlas dentro del scope activo (regla de scope guard del paso 5).
   - Si el linter reporta violations en archivos **tocados por la task activa**, son parte del scope y deben resolverse antes de `transition complete`.
2. La instrucción aplica a cualquier lenguaje; mencionar ejemplos representativos sin ser exhaustivo.

**Out**:
- No prescribir qué linter usar — inferir del repo.
- No agregar el linter run como paso universal para tasks que no tocan código.
- No modificar la regla de scope guard existente (paso 5 prohíbe trabajo fuera del spec).

## Criterios de Aceptación

- `loop-subcommand.md` incluye instrucción de lint local en paso 5 o 6 para tasks que modifican código fuente.
- La instrucción distingue violations en archivos propios de la task vs. archivos externos.
- Para violations externas, indica documentarlas como candidata a nueva task (no implementar).
- La instrucción menciona al menos dos ejemplos de linters de lenguajes distintos.
