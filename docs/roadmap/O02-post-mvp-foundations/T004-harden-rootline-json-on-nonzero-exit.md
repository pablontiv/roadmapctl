---
estado: Completed
tipo: task
---
# T004: Parsear JSON Rootline con exit non-zero

**Outcome**: [O02 Fundaciones post-MVP](README.md)
**Contribuye a**: CE3, INV2, INV3

## Preserva

- INV1: Errores Rootline siguen produciendo exit codes correctos en `roadmapctl`.
  - Verificar: tests de diagnostics/exit aggregation.

## Contexto

`internal/rootlinecli/client.go` hoy parsea JSON solo si el proceso Rootline sale con código cero. Rootline puede emitir JSON útil para validaciones fallidas; descartarlo degrada diagnostics y obliga a mensajes genéricos.

## Alcance

**In**:
1. Capturar stdout JSON aunque Rootline salga non-zero.
2. Devolver `JSONResult` junto con error estructurado cuando sea posible.
3. Permitir que `roadmapctl` emita diagnostics específicos desde JSON y preserve stderr/exit info.
4. Agregar tests con fake Rootline que retorna non-zero + JSON válido.

**Out**:
- Cambiar el contrato público de Rootline.
- Ocultar errores de ejecución reales sin JSON.

## Estado inicial esperado

- `runJSON` llama `c.run` y retorna error antes de `json.Unmarshal`.

## Criterios de Aceptación

- Un `rootline validate --all` inválido con JSON produce diagnostics parseados.
- Non-zero sin JSON sigue produciendo diagnostic de operación Rootline.
- JSON mode de `roadmapctl` sigue siendo un único objeto stdout.

## Fuente de verdad

- `internal/rootlinecli/client.go`
- `internal/rootlinecli/client_test.go`
- `internal/roadmap/dependencies.go`
