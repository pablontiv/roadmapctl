---
estado: Specified
tipo: task
---
# T035: Integrate — mapeo determinístico de `<type>` para commit_style conventional

**Contribuye a**: La derivación del tipo de commit conventional no depende de que la task tenga un prefijo (`feat:`, `fix:`, etc.); usa una tabla determinística basada en el contenido de la task.

## Contexto

`.claude/skills/integrate/SKILL.md` dice hoy: "el `type` se infiere del prefijo de la task (`feat`, `fix`, `docs`, `chore`, etc.)". Pero las tasks del repo rara vez llevan ese prefijo — sus títulos son descriptivos ("Add picokit dependency", "Migrate fsx to pathsec", "Loop force cache bypass"). La regla se queda sin entrada y el modelo improvisa.

En la sesión O25 inferí types desde el contenido: T001 (go get) → `chore`; T026-T029 (editan skill md) → `docs`; T002-T004 (mueven paquetes) → `refactor`. Funcionó porque hubo criterio del modelo; otro modelo podría inferir distinto y romper la consistencia del historial.

## Alcance

**In**:
1. En `.claude/skills/integrate/SKILL.md`, Fase 3 sección de derivación de mensaje (`commit_style: conventional`):
   - Reemplazar "inferir del prefijo de la task" por una tabla de mapeo determinístico desde el contenido de la task.
   - La tabla debe cubrir al menos estas filas:

     | Si la task… | `<type>` |
     |---|---|
     | Solo toca `.md`/docs/skills, sin código ejecutable | `docs` |
     | Agrega, bumpea o quita dependencias sin código fuente | `chore` |
     | Mueve, renombra o reemplaza paquetes preservando la API pública | `refactor` |
     | Agrega capability ejecutable nueva | `feat` |
     | Corrige un bug puntual sin cambiar API | `fix` |
     | No matchea ninguna categoría | `chore` (fallback con warning visible) |

   - Mantener la posibilidad de override explícito vía `commit_message` (ya soportado).
   - Mantener `<scope-corto>` derivado del Outcome (`O24`, `direct`, etc.).

**Out**:
- No agregar nuevas categorías al estándar conventional.
- No prescribir un linter de commit messages (eso es del repo).
- No cambiar el formato general del mensaje.

## Criterios de Aceptación

- Fase 3 contiene la tabla con las 6 filas mínimas listadas.
- El fallback `chore` con warning visible está explícitamente documentado.
- El override por `commit_message` queda preservado y mencionado.
- La frase "inferir del prefijo de la task" desaparece del skill.

## Fuente de verdad

- `.claude/skills/integrate/SKILL.md`
