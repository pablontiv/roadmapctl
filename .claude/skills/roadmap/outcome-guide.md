# Outcome Guide — Crear Objetivo de Roadmap

## Cuándo crear un Outcome

Crear un Outcome solo si agrupa una intención común y reduce complejidad:

- más de 5 tasks relacionadas,
- un resultado observable que requiere varias tasks,
- invariantes compartidas,
- necesidad de tracking visual por objetivo.

Si el trabajo cabe en 1–5 tasks auto-contenidas, crear tasks directas bajo `<roadmap-root>/`.

## Auto-numbering

El skill no calcula `OXX`. `roadmapctl materialize --dry-run` asigna el siguiente Outcome determinísticamente y muestra la ruta propuesta en `changes[]`; `--apply` crea el directorio/README si el dry-run fue aprobado.

## Estructura

```text
<roadmap-root>/
└── O01-nombre-del-objetivo/
    ├── README.md
    └── T001-task.md
```

## Template: Outcome README

```markdown
---
estado: Pending
tipo: outcome
---
# OXX: [Nombre del objetivo]

## Objetivo

[Resultado observable que existirá cuando todas las tasks estén completadas.]

## Criterios de Éxito

- CE1: [criterio verificable]
  - Verificar: [comando o procedimiento]
- CE2: [criterio verificable]
  - Verificar: [comando o procedimiento]

## Invariantes

- INV1: [propiedad que ninguna task debe romper]
  - Verificar: [comando o procedimiento]

## Alcance

**In**:
- [incluido]

**Out**:
- [excluido]

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-ejemplo.md) | Descripción breve |
```

## Reglas

- La tabla de tasks no incluye estado; el estado vive en el frontmatter de cada task.
- No crear subniveles bajo Outcome.
- Las dependencias entre tasks se declaran con `[[blocked_by:./TXXX-name.md]]` en la task bloqueada.
