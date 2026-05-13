# Modo Autónomo — Descomposición Simple

Cuando `$ARGUMENTS` no empieza con `pending|loop|plan`, generar una propuesta de roadmap con máximo dos niveles: Outcome opcional + Tasks.

## Paso 0: Bootstrap

Ejecutar primero el bootstrap de `SKILL.md`. Usar el JSON de `roadmapctl bootstrap` como fuente de configuración antes de cualquier análisis. No leer ni crear archivos de config legacy.

## Paso 1: Resolver intención

> **Razonamiento profundo**: Para proyectos complejos (>3 Outcomes anticipados o área técnica poco familiar), usar el `effort: xhigh` declarado en el frontmatter del skill. El effort se puede subir a `max` con `/effort max` si el scope es muy grande.

Determinar desde `$ARGUMENTS`:

- objetivo/capacidad a construir,
- repo target en workspace mode,
- documentación o código relevante,
- si el trabajo requiere Outcome o solo tasks.

En workspace mode, si el repo no es evidente, preguntar con opciones concretas.

## Paso 2: Absorber contexto acotado

Discovery determinista:

1. Reusar el JSON de `roadmapctl bootstrap` del repo target: `root`, `roadmap_root`, helpers, estados y opciones operacionales.
2. Leer `README*` de la raíz (máx. 3).
3. Buscar docs por keywords de `$ARGUMENTS` (máx. 8, preferir `docs/`, `research/`, `intent/`).
4. Leer Outcomes existentes relacionados bajo `<roadmap-root>/` (máx. 8) para evitar overlap.
5. Si el scope menciona código, leer manifests/entrypoints relevantes (máx. 10): `go.mod`, `package.json`, `Cargo.toml`, `pyproject.toml`, `justfile`, `Makefile`, `cmd/**`, `src/**`.
6. Si falta una decisión crítica, preguntar con opciones concretas.

No leer “todo el repo”.

## Paso 3: Normalizar vocabulario

El roadmap es implementación, no investigación. Traducir vocabulario exploratorio:

| Evitar | Usar |
|--------|------|
| hipótesis, premisa | requisito, objetivo |
| CAP-XX, LI-XX, H-XX | nombre técnico descriptivo |
| falsación, evidencia | verificación, criterio |
| fase/ciclo de investigación | eliminar |

Test: un desarrollador que no leyó la investigación entiende cada Outcome y Task.

## Paso 4: Elegir estructura

Reglas:

- 1–5 tasks auto-contenidas → tasks directas.
- Más de 5 tasks relacionadas → 1 Outcome + tasks.
- Objetivos independientes → varios Outcomes.
- Nunca crear niveles intermedios.

Leer [framework-reference.md](framework-reference.md) y aplicar sus criterios.

## Paso 5: Generar propuesta

Formato para tasks directas:

```text
TASKS DIRECTAS
├── T001: [task accionable] — [descripción 1 línea]
└── T002: [task accionable] — [descripción 1 línea]
```

Formato con Outcome:

```text
O01: [Objetivo observable]
├── Criterios de éxito:
│   ├── CE1: ... (verificar: ...)
│   └── CE2: ... (verificar: ...)
├── Invariantes:
│   └── INV1: ... (verificar: ...)
└── Tasks:
    ├── T001: [task accionable] — [descripción 1 línea]
    └── T002: [task accionable] — [descripción 1 línea]
```

Para cada task incluir:

- nombre,
- descripción de 1 línea,
- hard blockers `blocked_by` solo si existen y están justificados,
- criterios de aceptación principales.

No proponer `blocked_by` por orden sugerido, secuencia narrativa, relación temática o “conviene después de”. Para cada blocker propuesto, incluir la razón breve: qué fallaría objetivamente si la task se ejecuta antes.

## Paso 6: Validación antes de presentar

Verificar:

1. Cada task contribuye a un Outcome o resultado directo.
2. Cada criterio de éxito tiene al menos una task que lo implementa.
3. No hay tasks duplicadas.
4. Cada `blocked_by` propuesto es un hard blocker objetivo, no orden/contexto blando, y no forma ciclos evidentes.
5. Cada task cabe en una sesión.
6. La propuesta usa únicamente Outcomes y Tasks.

## Paso 7: Presentar para aprobación

Presentar propuesta fundamentada. No preguntar por la taxonomía; el framework ya define la granularidad.

Después de aprobación, informar que `/roadmap plan` materializa los archivos `.md`.
