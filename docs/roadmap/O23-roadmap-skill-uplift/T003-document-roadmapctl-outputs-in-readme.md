---
estado: Completed
tipo: task
---
# T003: Documentar outputs reales de roadmapctl en README + quitar `--repo .` redundante

**Outcome**: [O23 Roadmap skill: tool access, observability docs, README outputs](README.md)
**Contribuye a**: dar a usuarios humanos y a AI agents que leen el README una vista concreta del contrato real de cada comando (cómo se ve el output esperado), y eliminar ruido sintáctico (`--repo .` cuando el default del flag global ya es `.`).

## Preserva

- INV1: La estructura de secciones top-level del README (Quick Start, Core Idea, Layer Responsibilities, Installation, Auto-update, Commands, AI-Native, Skill Source, Development, Documentation, License) se preserva.
  - Verificar: `grep -nE "^## " README.md` antes y después produce la misma lista de headers en el mismo orden.
- INV2: Los comandos exactos (post-limpieza, sin `--repo .`) siguen siendo correctos cuando el cwd es la raíz del repo.
  - Verificar: ejecutar al menos uno de los comandos limpiados (e.g., `roadmapctl doctor`) desde la raíz del repo y confirmar exit 0 con output equivalente al output documentado.

## Contexto

`/home/shared/roadmapctl/README.md` actualmente muestra los comandos pero no su stdout. Usuarios y agentes que leen el README no tienen idea de qué formato esperar (texto tabular vs JSON, qué campos, qué exit codes); muchos terminan ejecutando el comando solo para descubrir el shape de output. Para un repo cuyo selling point es "stable JSON con versioned contracts" (línea ~67 del README), no mostrar los contratos en el README es una omisión grave.

Adicionalmente, todos los ejemplos del README cargan `--repo .` redundante. Confirmado vía `roadmapctl doctor --help`: el flag global `--repo` ya tiene default `"."`. La forma corta funciona sin diferencia desde la raíz del repo. El flag explícito es ruido que entrena a usuarios y agentes a tipear flags innecesarios.

Dos cambios acoplados en la misma task porque tocan el mismo archivo y los mismos bloques de ejemplos:

### (a) Quitar `--repo .` redundante

Eliminar en todos los lugares donde el ejemplo corre asumiendo cwd = raíz del repo. Conservar `--repo <path>` solo en ejemplos que demuestren explícitamente apuntar a otro repo (workspace mode, si existen).

### (b) Agregar bloques `# Output:`

Inmediatamente debajo de cada comando relevante en el mismo fenced-block group, agregar un bloque con el output real capturado. Comandos a cubrir (mínimo):

- `roadmapctl doctor` → output text default mostrando `status: ok`, errors/warnings counts (también se puede mostrar `--output json` brevemente si es útil).
- `roadmapctl check --strict` → text o JSON con `lint: total: 0`, `ok: true`.
- `roadmapctl pending` → tabla text default con 2-3 filas sample (ID, estado, titulo).
- `roadmapctl next` → text o JSON con top del `ready[]` queue.
- `roadmapctl decision` → output con secciones `ready` / `blocked_by_deps` / `blocked_by_estado`.
- `roadmapctl bootstrap --output json` → JSON trimmed con `root`, `roadmap_root`, `helpers`, status/config fields — alto valor porque es el "contrato" desde el que los agentes bootstrapean.
- `roadmapctl transition start ... --apply` → text o JSON con `from`, `to`, `applied: true`.

Cómo capturar los outputs: correr cada comando contra `/home/shared/roadmapctl/` mismo, copiar el output real (no redactar — el roadmap de este repo es público vía README). Cada bloque ≤15 líneas; truncar arrays largos con `…` si el output excede.

### Placement rule

El bloque output va inmediatamente después del comando en el mismo fenced-block group, prefixed con el comentario `# Output:` para que comando-y-resultado se vean inline cuando se scanea el README.

## Alcance

**In**:
1. Editar `/home/shared/roadmapctl/README.md` para eliminar `--repo .` redundante.
2. Capturar y agregar bloques `# Output:` para los 7 comandos enumerados, ubicados inline después de cada comando.
3. Verificar que el README sigue renderizando bien (markdownlint si está configurado, o smoke visual).

**Out**:
- Mover comandos entre secciones o restructurar el README.
- Agregar comandos nuevos al README (solo documentar los ya presentes).
- Editar otros archivos de documentación (`docs/cli-contract.md`, etc.).
- Actualizar `CHANGELOG.md` por esta task — entra en el changelog del próximo release agrupado, no por task individual.

## Estado inicial esperado

- `README.md` tiene múltiples ejemplos con `--repo .` (Quick Start ~líneas 36-58, Commands ~líneas 136-160).
- Ningún comando del README muestra output esperado.

## Criterios de Aceptación

- `grep -c -- "--repo \\." README.md` retorna 0, o únicamente cuenta apariciones en ejemplos que explícitamente demuestran workspace mode (apuntar a un repo distinto).
- Cada uno de los 7 comandos enumerados en Contexto tiene un bloque `# Output:` adyacente con output real capturado del repo `/home/shared/roadmapctl/`.
- Cada bloque output es ≤15 líneas (medible con `wc -l` sobre el rango del fenced block).
- `markdownlint README.md` exit 0 (si está configurado en el repo; si no, verificación visual via render en branch de PR).
- `diff` entre estado pre/post muestra solo cambios localizados (eliminación de `--repo .` y adición de output blocks), sin restructuración de secciones.

## Fuente de verdad

- `/home/shared/roadmapctl/README.md`
- `roadmapctl doctor --help` (referencia del default `"."` para `--repo`)
