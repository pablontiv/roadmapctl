# /roadmap plan

Materializa el plan de la conversación como archivos `.md` del roadmap. No implementa código.

Ruta normal autosuficiente: este archivo contiene el procedimiento operativo completo. No leer `common-logic.md` ni documentación de integración para ejecutar el flujo; esos documentos son referencia de mantenimiento/troubleshooting.

Materializar es una operación estructural. Está prohibido crear un único archivo
con una lista de tareas. Cada task debe tener su propio archivo `TXXX-*.md`.

## Fuente del plan

1. Contexto actual de conversación.

Si no hay plan, informar: "No hay plan en esta conversación. Primero planificar, luego ejecutar `/roadmap plan`." y parar.

## Workspace mode

Resolver repo target:

1. `--repo <name>` si fue dado.
2. Repo mencionado en el plan.
3. Si ambiguo, preguntar.

Usar `<abs-roadmap-root>` y `git -C <repo-path>`.

## Fase 1: Descomposición

1. Identificar el plan más reciente de la conversación.

2. Consultar numeración actual:
   ```bash
   rootline describe <roadmap-root> --field schema.id.next_by_pattern --output json
   ```
   Retorna `{"O*": "O14", "T*": "T014"}`.
   - Usar `O*` para el siguiente Outcome.
   - Usar `T*` como referencia inicial para tasks en outcomes nuevos.

   Para tasks en un **Outcome existente**:
   ```bash
   rootline describe <roadmap-root>/OXX-slug/ --field schema.id.next_by_pattern --output json
   ```
   Retorna `{"T*": "T009"}` — primer task disponible dentro de ese outcome.

   Para tasks en un **Outcome nuevo** (directorio aún no existe): comenzar tasks desde T001.

3. Aplicar `framework-reference.md`: máximo Outcome + Tasks por outcome.
4. Asignar slugs (kebab-case, sin prefijos O/T, sin `/` ni `..`) y numerar con los valores obtenidos.
5. Cada task: nombre, descripción, ACs principales, `hard_blockers` solo si hay dependencia objetiva real.

## Fase 2: Aprobación

Presentar árbol completo con números asignados + ACs:

```
O14-nombre-outcome/
├── README.md
├── T001-primera-task.md
│   - AC1: ...
└── T002-segunda-task.md
    - AC1: ...
```

**STOP obligatorio** con `AskUserQuestion` hasta aprobación explícita. No crear archivos antes.

## Fase 3: Materialización

**3.1 Re-confirmar numeración (antistaleness)**

```bash
rootline describe <roadmap-root> --field schema.id.next_by_pattern --output json
```

Si aparecieron nuevos archivos que cambian los números propuestos, informar al usuario y recalcular antes de continuar.

**3.2 Preflight obligatorio**

```bash
command -v roadmapctl
roadmapctl doctor --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si cualquier comando sale non-zero: detenerse, reportar exit code y diagnostics. No crear archivos.

**3.3 Escritura en paralelo**

Crear directorios padre si aplican, luego escribir con Write tool en paralelo:

- `OXX-slug/README.md`: frontmatter `tipo: outcome` + título + descripción/contexto (SIN `## Criterios de Aceptación` ni `## Tasks`). Ver template en `outcome-guide.md`.
- `OXX-slug/TXXX-slug.md`: frontmatter `estado: Specified` + título + descripción + `## Criterios de Aceptación` + contexto + scope + hard blockers si aplican. Ver template en `task-guide.md`.

**3.4 Validación por archivo**

```bash
rootline validate <path-creado>
```

Por cada archivo creado. Si falla: reportar y detener.

**3.5 Postcheck obligatorio**

```bash
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si falla: detenerse, reportar diagnostics. No commitear.

## Fase 4: Commit

```bash
git -C <repo-path> add <archivos .md creados>
git -C <repo-path> commit -m "chore(roadmap): create planning docs"
```

STOP. Informar: "Archivos de planificación creados. Ejecutar `/roadmap loop` cuando esté listo para implementar."
