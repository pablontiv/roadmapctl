# Outcome Guide — Crear Objetivo de Roadmap

## Cuándo crear un Outcome

Crear un Outcome solo si agrupa una intención común y reduce complejidad:

- más de 5 tasks relacionadas,
- un resultado observable que requiere varias tasks,
- invariantes compartidas,
- necesidad de tracking visual por objetivo.

Si el trabajo cabe en 1–5 tasks auto-contenidas, crear tasks directas bajo `<roadmap-root>/`.

## Auto-numbering

El skill no calcula `OXX`. `roadmapctl path-planning` (cuando esté disponible) asigna el siguiente Outcome determinísticamente y muestra la ruta propuesta. Mientras tanto, usar `roadmapctl next` como referencia para numbering determinístico.

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

[Contexto adicional, background, o motivación si aplica — prose libre sin secciones fijas.]
```

## Reglas

- El title del README usa el título del plan (`# [Nombre del objetivo]`); el identificador `OXX` vive en la ruta asignada por path-planning.
- El cuerpo materializado contiene descripción y contexto; no incluye `## Criterios de Aceptación` ni `## Tasks` (son vistas calculadas derivadas de los archivos `TXXX-*.md` hijas).
- El estado de `README.md` (outcome/index) no debe escribirse manualmente: se deriva desde las `TXXX-*` hijas y/o el estado del índice jerárquico.
- No crear subniveles bajo Outcome.
- Las dependencias duras entre tasks se declaran con `[[blocked_by:./TXXX-name.md]]` en la task bloqueada solo si la task no puede ejecutarse/validarse antes; orden sugerido o relación temática van en contexto/prose, no en `blocked_by`.
