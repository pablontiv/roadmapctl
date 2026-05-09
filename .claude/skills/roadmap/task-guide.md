# Task Guide — Crear Task AI-Ready

## Workflow

Este archivo define el contenido AI-ready de una task. No es el procedimiento primario para escribir archivos: `/roadmap plan` debe serializar el árbol aprobado a JSON y delegar rutas, numbering, creación de archivos, tabla del README y links `blocked_by` en `roadmapctl materialize`.

### Paso 1: Parsear argumentos

Extraer para el plan estructurado:

- **task-name**: slug kebab-case, ej. `add-k8s-phase`.
- **descripción**: qué debe hacer el agente.
- **outcome opcional**: slug/título del Outcome si la task pertenece a uno.

### Paso 2: Poblar plan estructurado

Para cada task, completar campos requeridos del schema `roadmapctl/materialize-plan`: slug, título, descripción, `preserves`, contexto, `scope_in`, `scope_out`, estado inicial, criterios de aceptación, fuentes de verdad y dependencias.

### Paso 3: Delegar materialización

`roadmapctl materialize --dry-run` asigna `TXXX`, detecta colisiones/escapes y muestra rutas. `roadmapctl materialize --apply` crea el archivo y actualiza el README padre si corresponde. El skill no debe ejecutar `rootline new` ni editar tablas manualmente.

## Dependencias

Usar `blocked_by` en la task bloqueada, siempre con path relativo explícito:

```markdown
[[blocked_by:./T001-prerequisite.md]]
```

Si la dependencia está en otro Outcome:

```markdown
[[blocked_by:../O01-setup/T001-prerequisite.md]]
```

Significa: la task actual no puede ejecutarse hasta que esa task esté completada.

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

[[blocked_by:./TXXX-prerequisite.md]] <!-- omitir si no hay dependencia -->

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
4. ¿Declara dependencias con `blocked_by` y path relativo explícito?
5. ¿Lista fuentes de verdad?
6. ¿Preserva invariantes relevantes?
