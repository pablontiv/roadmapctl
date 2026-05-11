# Task Guide — Crear Task AI-Ready

## Workflow

Este archivo define el contenido AI-ready de una task. El procedimiento primario (`/roadmap plan`) descompone el árbol en Outcomes/Tasks, pide aprobación del usuario, y luego **el skill escribe directamente los archivos `.md`** tras aprobación.

### Paso 1: Parsear argumentos

Extraer para el plan estructurado:

- **task-name**: slug kebab-case, ej. `add-k8s-phase`.
- **descripción**: qué debe hacer el agente.
- **outcome opcional**: slug/título del Outcome si la task pertenece a uno.

### Paso 2: Poblar plan estructurado

Para cada task, completar campos requeridos: slug, título, descripción, preserves, contexto, scope_in, scope_out, estado inicial, criterios de aceptación, fuentes de verdad y hard blockers opcionales.

### Paso 3: Escritura directa tras aprobación

Después de que el usuario aprueba explícitamente el árbol visual propuesto:

1. El skill escribe archivos `.md` directamente usando Write tool.
2. Cada Task crea un archivo `TXXX-task-slug.md` con template definido abajo.
3. El skill puede escribir archivos en paralelo si los parents (Outcomes) ya existen o fueron creados en el mismo batch.
4. Tras escribir, ejecutar `rootline validate <path>` sobre cada archivo crítico.
5. Ejecutar `roadmapctl check --strict` tras escribir todos los archivos para postcheck obligatorio.

## Dependencias duras

`blocked_by` significa hard blocker: la task actual **no debe ejecutarse** hasta que la task objetivo esté completada según `done_statuses`.

Antes de declarar cualquier `blocked_by`, responder:

```text
¿Qué fallaría objetivamente si ejecuto esta task antes?
```

Usar `blocked_by` solo si hay una respuesta concreta: falta una API, contrato, archivo, migración, test base o decisión sin la cual la task actual no puede validarse. Si la relación es orden sugerido, contexto, tema compartido, provenance, "conviene después de" o "usar su output si existe", no usar `blocked_by`; ponerlo en `Contexto`, `Fuente de verdad` o prose.

Si existe hard blocker, declararlo en la task bloqueada con path relativo explícito:

```markdown
[[blocked_by:./T001-prerequisite.md]]
```

Entre Outcomes:

```markdown
[[blocked_by:../O01-setup/T001-prerequisite.md]]
```

No usar targets bare como `[[blocked_by:T001-prerequisite]]`: rootline solo puede resolverlos por basename único y se rompen si hay duplicados.

## Template: Task File

```markdown
---
estado: Specified
tipo: task
---
# TXXX: [Descripción accionable]

**Outcome**: [OXX Nombre](README.md) <!-- omitir si es task directa -->
**Contribuye a**: [criterio de éxito del Outcome o resultado directo esperado]

[[blocked_by:./TXXX-prerequisite.md]] <!-- omitir salvo hard blocker objetivo -->

## Preserva

- INV1: [invariante a mantener]
  - Verificar: [comando o procedimiento]

## Contexto

[Contexto suficiente para que un agente ejecute esta task leyendo solo este archivo.]

## Alcance

**In**:
1. [acción concreta]
2. [acción concreta]

**Out**:
- [límite explícito]

## Estado inicial esperado

- [precondición observable]

## Criterios de Aceptación

- [AC binario con comando/check esperado]
- [AC binario con comando/check esperado]

## Fuente de verdad

- [paths a leer/modificar]
```

## Estados

| Estado | Cuándo |
|--------|--------|
| Pending | Task creada, aún no especificada completamente |
| Specified | Lista para implementar |
| In Progress | Ejecución en curso |
| Completed | Ejecutada y verificada |
| Blocked | Bloqueada por dependencia o condición externa |
| On Hold | Diferida intencionalmente |
| Obsolete | Ya no aplica |

## Checklist

Antes de finalizar una task, verificar:

1. ¿Cabe en una sesión?
2. ¿Contiene todo el contexto?
3. ¿Los ACs son pass/fail?
4. Si declara `blocked_by`, ¿hay respuesta concreta a "qué fallaría objetivamente si ejecuto esta task antes"?
5. ¿Los links `blocked_by` son paths relativos explícitos y no orden/contexto blando?
6. ¿Lista fuentes de verdad?
7. ¿Preserva invariantes relevantes?
