---
estado: Completed
tipo: task
---
# T006: Revisar capacidades genéricas necesarias en Rootline

**Outcome**: [Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill](README.md)

## Preserva

- Rootline no incorpora comandos ni conceptos específicos de roadmapctl.

## Contexto

El nuevo flujo depende de Rootline como motor genérico para schema/frontmatter/sections/links y de roadmapctl como wrapper de dominio.

## Alcance

**In**:
1. Verificar rootline new con match patterns relevantes.
2. Verificar rootline set para secciones schema-backed.
3. Aclarar o corregir --create si no crea archivos como documentado.
4. Probar validate/query/tree/graph para las necesidades del roadmap.

**Out**:
1. No agregar lógica de Outcome/Task a Rootline.

## Estado inicial esperado

Smoke tests mostraron gaps o comportamientos a aclarar en rootline set --create y secciones.

## Criterios de Aceptación

- Hay evidencia/test de Rootline para frontmatter requerido por .stem.
- Hay evidencia/test de Rootline para secciones requeridas y set de secciones.
- La documentación de Rootline coincide con el comportamiento de --create.
- Cualquier gap se documenta como prerequisite o se corrige genéricamente en Rootline.

## Fuente de verdad

- /home/shared/rootline/docs/new.md
- /home/shared/rootline/docs/set.md
- /home/shared/rootline/docs/validate.md
- /home/shared/rootline/cmd/rootline/new.go
- /home/shared/rootline/cmd/rootline/set.go
