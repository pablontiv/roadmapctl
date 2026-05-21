# /roadmap plan

> Referencias obligatorias antes de escribir: [bootstrap-reference.md](bootstrap-reference.md)

## Fase 1: Bootstrap y detección gitflow

Ejecutar bootstrap obligatoriamente para detectar configuración y diagnostics:

```bash
roadmapctl bootstrap --repo <repo-path> --output json
```

Si bootstrap devuelve diagnóstico `RMC_GITFLOW_NOT_CONFIGURED`: ejecutar wizard de adopción gitflow (ver [bootstrap-reference.md](bootstrap-reference.md) sección "Wizard de adopción gitflow") → re-ejecutar bootstrap → continuar con el flujo normal del subcomando.

---

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
roadmapctl doctor --repo <repo-path> --output json --strict
roadmapctl check --repo <repo-path> --output json --strict
```

Si cualquier comando sale non-zero: detenerse, reportar exit code y diagnostics. No crear archivos.

> **Stem legacy:** Si `doctor`/`check` falla con `RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED` o
> `RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY`, ejecutar:
> ```bash
> roadmapctl bootstrap --repo <repo-path> --yes
> ```
> y reintentar el preflight. En modo autónomo (`--yes`) la reparación aplica sin prompt.
> Si el `.stem` tiene campos custom no reconocidos, bootstrap emite `RMC_BOOTSTRAP_REPAIR_UNSUPPORTED_STEM`
> y no modifica nada — escalar al usuario en ese caso.

**3.3 Escritura en paralelo — Model de Dispatch de Agents**

Ejecutar la materialización del plan mediante un modelo de waves + Agent dispatch:

**Wave 0: Materialización de Outcomes (Coordinador)**

- El coordinador escribe los `OXX-slug/README.md` directamente.
- Alternativa: si hay múltiples Outcomes con gran volumen, despatchar 1 Agent para materializarlos todos.
- Los Outcomes deben existir antes de crear Tasks (precondición estructural).
- Template: frontmatter `tipo: outcome` + título + descripción/contexto (SIN `## Criterios de Aceptación` ni `## Tasks`). Ver `outcome-guide.md`.

**Wave 1+: Materialización de Tasks (Agents en paralelo)**

- Particionar el conjunto de Tasks en subsets: ~3-5 archivos por Agent, o 1 Agent por Outcome + sus Tasks asociadas.
- Despatchar N Agents en paralelo, cada uno responsable por un subset disjunto de Task files.
- Cada Agent recibe en su prompt:
  - Lista de `(path, contenido_completo)` — el contenido está completamente decidido por el coordinador.
  - No recalculan contenido; materializan exactamente como se especifica.
  - Template de prompt:

```
Eres un agente de materialización del roadmap. Escribe los siguientes archivos exactamente como se indica.

Para cada archivo:
  - Path: <path relativo al repo>
  - Contenido: <contenido completo del archivo>

Instrucciones:
- Escribe cada archivo con la herramienta Write. Usa paths absolutos: /home/shared/<repo>/<path>
- Confirma cada archivo escrito con "✓ <path>"
- No interpretes ni modifiques el contenido
- Si un archivo ya existe con contenido idéntico, escríbelo igualmente
- No hagas nada más: solo escribe los archivos listados

Archivos a escribir:
[LISTA DE ARCHIVOS CON CONTENIDO]
```

- Cada Agent hace **1 Write call por archivo** de su subset.
- El coordinador espera a que todos los Agents terminen.

**Post-Wave: Validación y Postcheck global (Coordinador)**

- Ver secciones 3.4 y 3.5.

**3.4 Validación por archivo**

```bash
rootline validate <path-creado>
```

Por cada archivo creado. Si falla: reportar y detener.

**3.5 Postcheck obligatorio**

```bash
roadmapctl check --repo <repo-path> --output json --strict
```

Si falla: detenerse, reportar diagnostics. No commitear.

## Fase 4: Commit

```bash
git -C <repo-path> add <archivos .md creados>
git -C <repo-path> commit -m "chore(roadmap): create planning docs"
```

STOP. Informar: "Archivos de planificación creados. Ejecutar `/roadmap loop` cuando esté listo para implementar."

### Prohibición post-plan

La aprobación del árbol propuesto autoriza **solo** la creación de archivos `.md`. Está explícitamente prohibido:

- Llamar `roadmapctl transition start` o `roadmapctl transition complete` sobre cualquier task
- Modificar el contenido de los archivos descritos por las tasks (código, documentación, skills)
- Continuar implementando cualquier lógica, código o documentación de las tasks recién creadas
- Tratar la aprobación del árbol como autorización para implementar

Todo lo que sigue a la creación de archivos `.md` es responsabilidad exclusiva de `/roadmap loop`.
