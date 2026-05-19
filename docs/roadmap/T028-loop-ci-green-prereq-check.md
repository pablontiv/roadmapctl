---
estado: Specified
tipo: task
---
# T028: Loop — verificar CI verde antes de ejecutar tasks con invariante de CI

**Contribuye a**: El loop no intenta tagear/publicar/desplegar cuando CI tiene failures pre-existentes no relacionados con la task activa.

## Contexto

Cuando una task tiene un INV o AC que requiere "CI verde en el branch" (e.g. "tag apunta a commit con CI verde", "deploy requiere pipeline verde"), el loop asume implícitamente que CI pasó. Si hay failures pre-existentes en el branch que no son responsabilidad de la task activa, el loop los descubrirá recién después de ejecutar la task y hacer push — desperdiciando trabajo y produciendo commits de fix no planificados.

Ejemplo real: T003 de picokit tenía INV1 "tag apunta a commit con CI verde". El branch tenía failures de lint y go.sum desde commits anteriores que nunca tuvieron un run CI verde en código Go. El loop ejecutó T003 asumiendo CI verde y descubrió los failures solo al intentar tagear.

## Alcance

**In**:
1. En `loop-subcommand.md`, Fase 1 Discovery, agregar un paso de pre-check selectivo:
   - Al procesar `ready[]`, si una task contiene en su sección `## Preserva` o `## Criterios de Aceptación` las palabras "CI verde", "CI pasa", "pipeline verde" o equivalente, ejecutar:
     ```bash
     gh run list --branch <base_branch> --limit 5 --json status,conclusion,headSha
     ```
   - Si ningún run reciente tiene `conclusion: "success"`, reportar la condición como bloqueante informativo y sugerir resolver CI antes de ejecutar la task.
   - Si `gh` no está disponible, advertir y continuar (no bloquear en entornos sin GitHub CLI).
2. El pre-check es selectivo (solo tasks con la condición explícita), no universal.

**Out**:
- No agregar el check a todas las tasks — solo las que declaran CI como invariante.
- No bloquear automáticamente; reportar y dejar que `autonomy` decida.

## Criterios de Aceptación

- `loop-subcommand.md` Fase 1 incluye el pre-check selectivo de CI verde.
- La condición de detección está definida (palabras clave en `## Preserva` / `## Criterios de Aceptación`).
- El comportamiento cuando `gh` no está disponible está especificado (advertir, no bloquear).
- La regla no requiere CI verde universal — solo para tasks que lo declaran explícitamente.
