---
tipo: task
estado: Specified
---

# T014 Fix loop parallel waves: Agent calls explícitas

Reescribir la sección "Parallel waves" de `loop-subcommand.md` para que el mecanismo de ejecución paralela sea explícito: llamadas paralelas al tool `Agent`, una por task de la wave.

## Criterios de Aceptación

- AC1: La sección "Parallel waves" especifica explícitamente que cada wave se despacha con llamadas paralelas al tool `Agent`, una por task
- AC2: La condición de worktrees se elimina como requisito de entrada y pasa a ser señal de dependencia faltante (no intentar merge)
- AC3: `parallel == false` mantiene comportamiento secuencial intacto
- AC4: `./scripts/sync-roadmap-skill.sh --install && ./scripts/sync-roadmap-skill.sh --check` pasan sin error
- AC5: Headless pi test definido en `SKILL.md` pasa (bootstrap + preflight)

## Especificación Técnica

Archivo: `.claude/skills/roadmap/loop-subcommand.md`

Reemplazar el bloque actual de "Parallel waves":

```
- La ejecución paralela debe usar worktrees aislados o una ruta de integración
  equivalente con control de conflictos. Si no hay aislamiento/control seguro,
  ejecutar secuencialmente aunque `parallel == true`.
```

Por:

```
Si `parallel == true`, ejecutar cada wave despachando llamadas paralelas al tool
`Agent` — una por task de la wave. Las tasks de una wave son independientes por
definición (`roadmapctl next` garantiza ausencia de `blocked_by` entre ellas),
por lo que Agent calls paralelas sobre archivos distintos son la ruta correcta
sin necesidad de worktrees.

Si dos tasks de la misma wave producen conflicto al integrar, tratar como
dependencia faltante según el modo de autonomía — no usar worktrees para forzar
el merge.
```

Después del cambio: `sync-roadmap-skill.sh --install`.
