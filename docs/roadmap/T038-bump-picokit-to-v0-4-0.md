---
estado: Completed
tipo: task
---
# T038: Bump picokit de v0.2.0 a v0.4.0

**Contribuye a**: cerrar el desfase del consumidor con la librería; aprovechar coverage-spec v1.1 (auto-discovery) y la firma variadic de `autoupdate.New`.

**Dependencia externa** (no expresable como `blocked_by` por estar en otro repo): requiere que `picokit/O04-autoupdate-envdisable-optional-and-windows-fix/T003-release-v0-4-0` esté Completed y el tag `v0.4.0` publicado en GitHub.

## Contexto

`/home/shared/roadmapctl/go.mod` apunta a `github.com/pablontiv/picokit v0.2.0` mientras que picokit ya publicó v0.2.1 y v0.3.0. La task espera el tag v0.4.0 (release que produce O04 de picokit con fix Windows + envDisable opcional) para bumpear de una sola vez al último.

El call-site `internal/cli/cli.go:5431` actual pasa los tres args (`"pablontiv/roadmapctl", "roadmapctl", "ROADMAPCTL_NO_UPDATE"`); seguirá compilando contra v0.4.0 sin modificación (la firma variadic acepta args explícitos).

Beneficios del bump:
- `pkcov` con coverage-spec v1.1 (auto-discovery) — puede simplificar `.coverage-floors.toml` si se aprovecha.
- Fix Windows ya incluido (no impacta directamente porque roadmapctl no se distribuye para Windows hoy, pero asegura consistencia futura).

## Alcance

**In**:

1. `cd /home/shared/roadmapctl && go get github.com/pablontiv/picokit@v0.4.0`.
2. `go mod tidy`.
3. Correr suite local:
   - `just check`
   - `just test`
   - `just coverage-check` (`pkcov` ya activo en O29; v0.4.0 trae auto-discovery v1.1).
4. Si auto-discovery cambia la salida de `pkcov`, revisar `.coverage-floors.toml`: opcionalmente simplificar entradas que pkcov puede inferir. No es obligatorio en esta task — basta con que el gate siga verde.
5. Push del commit con mensaje `chore(deps): bump picokit to v0.4.0`.

**Out**:
- No tocar el wiring de autoupdate ni cambiar `ROADMAPCTL_NO_UPDATE` — la firma variadic mantiene compatibilidad.
- No refactorizar `.coverage-floors.toml` más allá de lo necesario.

## Estado inicial esperado

- `go.mod` declara `github.com/pablontiv/picokit v0.2.0`.
- Tag picokit `v0.4.0` publicado (precondición: O04 cerrado).

## Criterios de Aceptación

- `go.mod` declara `github.com/pablontiv/picokit v0.4.0`.
- `go mod tidy` no introduce líneas extra.
- `just check && just test && just coverage-check` exit 0.
- CI verde tras push.

## Fuente de verdad

- `/home/shared/roadmapctl/go.mod`
- `/home/shared/roadmapctl/go.sum`
- `/home/shared/roadmapctl/.coverage-floors.toml` (revisar tras bump)
