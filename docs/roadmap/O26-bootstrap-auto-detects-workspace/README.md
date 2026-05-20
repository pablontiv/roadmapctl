---
tipo: outcome
---
# O26: Bootstrap auto-detects workspace mode

`roadmapctl bootstrap` (sin subcomando y sin flag `--workspace`) debe reconocer automáticamente cuando se lo invoca en la raíz de un workspace multi-repo, en lugar de fallar con `RMC_CONFIG_MISSING`.

## Motivación

Hoy `buildBootstrapConfig` (`internal/cli/bootstrap.go:244`) llama directo a `config.Load(repo)` sin consultar `options.Workspace` ni inspeccionar la estructura del directorio. Cuando alguien parado en un directorio raíz sin `.git/` pero con varios members configurados (cada uno con su `.git/` + `docs/roadmap/.roadmapctl.toml`) ejecuta `roadmapctl bootstrap`, la CLI falla con `RMC_CONFIG_MISSING` apuntando a un archivo workspace-level que no debe existir.

Caso real: `/home/shared/` agrupa ~25 repos hermanos. Siete ya tienen roadmap propio (`backscroll`, `cartyx`, `crossbeam`, `picokit`, `pinata`, `roadmapctl`, `rootline`). Correr `/roadmap` desde el shared root falla, aunque la información necesaria para resolver workspace existe en disco.

## Diseño

La detección es una propiedad emergente de la estructura del filesystem, no un archivo TOML adicional:

```
IsWorkspaceRoot(root) := !exists(root/.git)
                        AND existe ≥1 subdir D con D/.git
                            AND D/docs/roadmap/.roadmapctl.toml
```

Las funciones de scan (`workspaceRepoRoots` en `pending.go:114`) ya existen para el comando `pending --workspace`; el cambio es extraerlas a un paquete compartido (`internal/workspace/`) y consumirlas también desde bootstrap.

## Alcance

In:
- Refactor estructural que mueve la detección de workspace a un paquete reutilizable.
- Nueva rama workspace en `buildBootstrapConfig`, activada por `options.Workspace` o por auto-detección.
- Diagnóstico nuevo `RMC_WORKSPACE_EMPTY` cuando se detecta raíz workspace pero ningún member tiene config válido.
- Cobertura de tests para los nuevos paths.

Out:
- No se introduce ningún archivo TOML workspace-level. La detección no debe depender de él.
- No se modifica el contrato del skill `/roadmap`, ni `pending-subcommand.md`, ni otros subcomandos.
- No se cambia el comportamiento single-repo existente; los tests actuales deben seguir verdes sin tocarlos.
