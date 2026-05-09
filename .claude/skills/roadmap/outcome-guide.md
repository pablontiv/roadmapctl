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
tipo: outcome
---
# [Nombre del objetivo]

[Descripción del resultado observable que existirá cuando todas las tasks estén completadas.]

## Criterios de Aceptación

- [criterio verificable]

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-ejemplo.md) | Descripción breve |
```

## Reglas

- La tabla de tasks no incluye estado; el estado vive en el frontmatter de cada task.
- El título del README usa el título del plan (`# [Nombre del objetivo]`); el identificador `OXX` vive en la ruta asignada por `roadmapctl materialize`.
- El cuerpo materializado contiene descripción, `## Criterios de Aceptación` y `## Tasks`; no agregar secciones prose-only al template si el renderer no las emite.
- El estado de `README.md` (outcome/index) no debe escribirse manualmente: se deriva desde las `TXXX-*` hijas y/o el estado del índice jerárquico.
- No crear subniveles bajo Outcome.
- Las dependencias entre tasks se declaran con `[[blocked_by:./TXXX-name.md]]` en la task bloqueada.
