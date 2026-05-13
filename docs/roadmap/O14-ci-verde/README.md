---
tipo: outcome
---
# O14: CI verde

Todos los jobs del pipeline de CI pasan en master y el job `release` puede correr,
generando el primer GitHub Release con binarios, checksums y attestations.

Tres bloqueos impiden actualmente que CI esté verde: 28 issues detectados por
golangci-lint (linter ahora activo por primera vez), ~20 tests que fallan porque
el job `ci / Test` de crossbeam no instala rootline, y cobertura en 84.4% por
debajo del umbral de 85%.
