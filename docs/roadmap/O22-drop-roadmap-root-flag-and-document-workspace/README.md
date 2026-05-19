---
tipo: outcome
---
# O22: Drop `--roadmap-root` flag and document workspace multi-repo

Investigación con backscroll de sesiones anteriores confirmó cero casos reales de uso del flag `--roadmap-root` con valor distinto al default `docs/roadmap`. Todas las invocaciones del skill lo pasan redundantemente con el default. El flag es plumbing inútil.

Adicionalmente, los dos GitHub issues abiertos (#2 multi-repo workspace commit routing, #3 `.roadmapctl.toml` en repo root) no requieren features nuevas: el modelo correcto es que cada repo del workspace tenga su propio roadmap completo bajo `<repo>/docs/roadmap/`. Esta convención debe documentarse explícitamente.

El resultado observable cuando todas las tasks estén completadas:

- `roadmapctl --help` no muestra `--roadmap-root` y `roadmapctl <cmd> --roadmap-root <path>` falla con "unknown flag"
- `.claude/skills/roadmap/*.md` no pasa `--roadmap-root` en ningún ejemplo
- README contiene sección "Workspace mode" documentando que cada repo participante mantiene su propio roadmap
- SKILL.md Paso 0 aclara que el loop se invoca por repo en workspace mode
- `docs/cli-contract.md` no menciona el flag y registra `docs/roadmap/` como convención fija
- GitHub issues #2 y #3 cerrados con comentarios explicativos

Sin cambios al modelo workspace existente: `workspaceRepoRoots()` (`internal/cli/pending.go`) ya descubre repos hermanos vía `.git`; cada repo carga su propia config; el loop opera repo-by-repo. Solo se documenta el comportamiento ya existente y se elimina el flag noise.
