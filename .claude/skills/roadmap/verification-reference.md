> Cargar solo si vas a modificar el skill `/roadmap` o sus guards `roadmapctl`.

# Verification Reference

## Verificación obligatoria al modificar este skill

Todo cambio al skill `/roadmap` o a sus guards `roadmapctl` debe probarse con Pi headless antes de commit/release. No alcanza con grep.

Ejecutar desde el repo canónico:

```bash
./scripts/sync-roadmap-skill.sh --install
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/roadmap/SKILL.md --tools read,bash -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: the user asks "loop autonomo" in this repository. Do not modify files and do not run git commit/push. Perform only the bootstrap and the required preflight checks from the skill, then stop. In your final answer, list the exact commands you ran and whether roadmapctl doctor/check were required and passed.'
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/roadmap/SKILL.md --tools read,bash -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: there is an already approved plan to materialize one direct task, and the user says "crea las tareas". Do not create or modify files and do not run git commit/push. Perform only bootstrap and the required preflight checks that must happen before any roadmap write, then stop. In your final answer, list exact commands run and whether roadmapctl doctor/check were required and passed.'
```

La evidencia debe mostrar que `roadmapctl doctor` y `roadmapctl check` fueron requeridos y pasaron antes de loop/materialización, sin modificar archivos.

Antes de hacer push, `golangci-lint run ./...` debe reportar 0 issues (CI lint gate).

El coverage gate usa `pkcov` con floors por paquete en `.coverage-floors.toml` (mínimo global: 85%):

```bash
go test ./... -coverprofile=coverage.out
go run github.com/pablontiv/picokit/cmd/pkcov check
```

O equivalentemente: `just coverage-check`.
