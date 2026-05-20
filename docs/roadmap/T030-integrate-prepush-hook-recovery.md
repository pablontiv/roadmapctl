---
estado: Completed
tipo: task
---
# T030: Integrate — recuperación cuando el pre-push hook rechaza el push

**Contribuye a**: El skill /integrate no silencia ni bypassa un pre-push hook que rechaza por cambios coordinados faltantes; intenta una recuperación estructurada o se detiene reportando el motivo literal.

## Contexto

Cuando un pre-push hook del proyecto rechaza el push exigiendo cambios coordinados (p. ej. "X cambió pero Y no fue actualizado"), hoy /integrate no tiene política definida: el modelo puede optar por `--no-verify` (incorrecto), por amend (peligroso) o por detenerse sin diagnóstico estructurado.

Ejemplo real (sesión O25): el push de T002 fue rechazado por el pre-push hook del repo roadmapctl porque `internal/` cambió sin un cambio coordinado en `.claude/skills/roadmap/`. La salida del hook indicaba el path requerido; bastaba un commit doc complementario para satisfacerlo. La recuperación se hizo a mano en el loop, no por el skill.

La política base es no-bypass: el hook es del proyecto, el skill respeta su decisión.

## Alcance

**In**:
1. En `.claude/skills/integrate/SKILL.md`, Fase 4 (Push):
   - Si `git push` sale con exit ≠ 0 y stderr contiene la palabra `hook` o `pre-push`, NO degradar a `--no-verify` ni a `--force`.
   - Intentar parsear stderr buscando paths o nombres de archivo mencionados como requeridos (heurística simple: líneas que matchean un path-like + verbo como `update`, `change`, `required`, `not updated`).
   - Si detecta paths candidatos, proponer al caller (vía diagnostic informativo) un commit complementario sobre esos paths dentro del mismo push range y permitir reintento.
   - Si no parsea nada estructurado, emitir diagnostic `INTEGRATE_HOOK_REJECTED` con el mensaje literal del hook y detenerse — el operador decide cómo proceder.
2. Agregar entrada en la tabla "Errores comunes" con ID `INTEGRATE_HOOK_REJECTED`, causa ("pre-push hook del proyecto exige cambios coordinados"), y acción recomendada.

**Out**:
- No modificar el comportamiento de retry para `RMC_INTEGRATE_PUSH_REJECTED` (que cubre divergencia con remote — ese sí permite rebase + reintento).
- No agregar bypass automático ni interactivo.
- No modificar otros skills.

## Criterios de Aceptación

- `.claude/skills/integrate/SKILL.md` Fase 4 documenta el manejo de rechazo por hook diferenciándolo del rechazo por divergencia.
- La tabla "Errores comunes" incluye `INTEGRATE_HOOK_REJECTED` con causa, síntoma y acción.
- El skill establece explícitamente: prohibido `--no-verify` o `--force-with-lease` automático.
- La recuperación por parseo de stderr describe la heurística (qué buscar en la salida) sin acoplarse a un hook específico.

## Fuente de verdad

- `.claude/skills/integrate/SKILL.md`
