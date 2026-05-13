---
tipo: outcome
---
# O15: CI smoke verde (macOS + Windows)

Los tres jobs de smoke CI pasan en todos los sistemas operativos: ubuntu, macOS y Windows.

Tres problemas bloquean los smoke jobs de macOS y Windows después de O14:
una regla de gosec inválida para golangci-lint v2.10.1 que rompe `ci / Lint`,
tests de colisión de case que fallan en filesystems case-insensitive (macOS),
y paths con backslash en JSON que fallan `AssertNoBackslashes` en Windows.
