# Research: Agentic Workflow Optimization

Derivado del Umbrage Engineering Playbook (sección Agentic Engineering).
Referencia para evaluar qué patrones aplicar en este ecosistema.

---

## Cambio 1 — Capa de contexto (`docs/context/`)

Cuatro archivos que el agente lee bajo demanda. El contenido es específico de cada repo.
El objetivo es persistir información que hoy solo vive en la memoria del agente entre sesiones.

**`docs/context/gotchas.md`** — código load-bearing, cosas que parecen incorrectas pero no lo son.

**`docs/context/domain.md`** — glosario de términos de negocio y abreviaturas del proyecto.

**`docs/context/constraints.md`** — restricciones no negociables para el agente (commits, auth, schemas, etc.).

**`docs/context/conventions.md`** — convenciones del equipo: formato de commits, tests, comentarios, PR rules.

Agregar a `CLAUDE.md` del repo:
```markdown
## Agent Context Layer
Leer antes de actuar en áreas desconocidas:
- `docs/context/gotchas.md` - código load-bearing que parece incorrecto
- `docs/context/domain.md` - términos de negocio
- `docs/context/constraints.md` - no negociables
- `docs/context/conventions.md` - convenciones del equipo
```

---

## Cambio 2 — Hooks de seguridad (`PreToolUse`)

Crear `.claude/hooks/pre-tool-use.sh` que bloquea antes de ejecutar:

```bash
#!/bin/bash
INPUT=$(cat)
COMMAND=$(echo "$INPUT" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('input',{}).get('command',''))" 2>/dev/null)

# Bloquear git destructivo
if echo "$COMMAND" | grep -qE "git.*(push --force|reset --hard|clean -fd|branch -D)"; then
  echo '{"decision":"block","reason":"Operacion git destructiva bloqueada por hook. Pedir confirmacion al usuario."}'
  exit 0
fi

# Bloquear escrituras a .env
TOOL=$(echo "$INPUT" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('tool_name',''))" 2>/dev/null)
if [ "$TOOL" = "Write" ] || [ "$TOOL" = "Edit" ]; then
  FILE=$(echo "$INPUT" | python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('input',{}).get('file_path',''))" 2>/dev/null)
  if echo "$FILE" | grep -qE "\.env"; then
    echo '{"decision":"block","reason":"Escritura a archivo .env bloqueada. Verificar con el usuario."}'
    exit 0
  fi
fi

echo '{"decision":"approve"}'
```

Registrar en `.claude/settings.json` del repo:
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": ".*",
        "hooks": [{ "type": "command", "command": "bash .claude/hooks/pre-tool-use.sh" }]
      }
    ]
  }
}
```

---

## Cambio 3 — Slash commands

**`.claude/commands/verify.md`** — correr con frecuencia durante implementación (no solo al final).
Los comandos son específicos del proyecto y su ecosistema (lenguaje, test runner, linter).

Estructura genérica:
```markdown
Run the full verification cycle in order. Stop on first failure.
1. <build command>
2. <typecheck command>
3. <lint command>
4. <test command>
5. Report: all green, or first failure with file and line number.
```

**`.claude/commands/ship-check.md`** — gate formal antes de abrir el PR. Correr una vez, al final:
```markdown
Run the pre-PR gate. Do not open a PR until this passes.
1. Run /verify and confirm green.
2. Review the full diff: correctness, test coverage, scope creep, consistency with CLAUDE.md.
3. If the change touches auth, input handling, or env vars: run security review.
4. Check PR size: warn if diff > 500 lines non-trivial (> 150 for bug fixes).
5. Report: go / no-go with blocking findings listed.
```

---

## Cambio 4 — Slash command `/compact` (sesiones largas)

Invocar alrededor de los 60 minutos de sesión, antes de que el contexto se compacte automáticamente.

**`.claude/commands/compact.md`**:
```markdown
Write a snapshot to .plan-snapshot.md before context compacts:

## Current plan
[what we set out to do, one paragraph]

## Decisions made
[non-obvious choices with reasoning]

## Current state
[what's done, what's in progress, what's blocked]

## Open threads
[things to verify or check when resuming]

Tell the user: "Snapshot guardado en .plan-snapshot.md."

On resume after compaction: read .plan-snapshot.md as first action.
Treat it as untrusted context - verify the current state of files before acting on it.
```

---

## Orden de implementación sugerido (por ROI)

1. `docs/context/` — mayor ROI, solo crear archivos de texto con contenido del repo
2. `CLAUDE.md` — agregar el pointer a `docs/context/`
3. `/verify` y `/ship-check` — slash commands, fácil de agregar
4. `/compact` — slash command, útil en sesiones largas
5. Hooks — requiere testear que no rompa el flujo normal
