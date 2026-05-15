---
estado: Specified
tipo: task
---
# T005: Add `--field` scalar extraction to `roadmapctl bootstrap`

**Outcome**: [O20 Field-agnostic integration](README.md)
**Contribuye a**: consumers pueden extraer un valor escalar de bootstrap sin depender de jq o python3

## Preserva

- INV1: `roadmapctl bootstrap --output json` sin `--field` retorna el mismo objeto completo de siempre
  - Verificar: `roadmapctl bootstrap --repo /home/shared/rootline --roadmap-root docs/roadmap --output json | python3 -c "import json,sys; r=json.load(sys.stdin); assert 'roadmap_root' in r"`
- INV2: `go test ./...` verde
  - Verificar: `cd /home/shared/roadmapctl && go test ./...`

## Contexto

`roadmapctl bootstrap --output json` retorna un objeto JSON con ~15 campos. Para extraer un valor escalar (e.g. `roadmap_root`, `helpers.where_leaf`) un consumer necesita `jq` o `python3`. rootline ya tiene `--field` que hace esto: `rootline describe <path> --field schema.id.next_by_pattern --output json`. roadmapctl debería seguir el mismo patrón.

Campos de uso frecuente en el skill:
- `roadmap_root`
- `root`
- `helpers.where_leaf`
- `helpers.where_not_done`
- `helpers.where_active`

## Alcance

**In**:
1. Agregar flag `--field <dot-path>` al comando `roadmapctl bootstrap`
2. Cuando `--field` está presente: extraer el valor en el path indicado del JSON de bootstrap y retornarlo como scalar (string o número) en stdout, sin wrapper JSON
3. Soportar paths simples (`roadmap_root`) y anidados con punto (`helpers.where_leaf`)
4. Si el campo no existe: exit 1 con mensaje claro en stderr
5. Si el valor es un objeto o array: exit 1 con mensaje claro (no serializar)
6. Tests unitarios y golden

**Out**:
- No cambiar el comportamiento de `--output json` sin `--field`
- No implementar path syntax con arrays (e.g. `diagnostics[0]`)
- No agregar `--field` a otros comandos (scope: solo bootstrap)

## Criterios de Aceptación

- `roadmapctl bootstrap --repo /home/shared/rootline --roadmap-root docs/roadmap --field roadmap_root` imprime `/home/shared/rootline/docs/roadmap` sin JSON wrapper
- `roadmapctl bootstrap --repo /home/shared/rootline --roadmap-root docs/roadmap --field helpers.where_leaf` imprime `isIndex == false`
- `roadmapctl bootstrap --repo /home/shared/rootline --roadmap-root docs/roadmap --field nonexistent` sale con exit 1
- `roadmapctl bootstrap --repo /home/shared/rootline --roadmap-root docs/roadmap --field diagnostics` sale con exit 1 (valor es array)
- `roadmapctl bootstrap --repo /home/shared/rootline --roadmap-root docs/roadmap --output json` (sin `--field`) retorna JSON completo sin cambios
- `go test ./internal/cli/...` verde

## Fuente de verdad

- `internal/cli/bootstrap.go` — comando bootstrap, flag `--field`, función de extracción
- `internal/cli/bootstrap_test.go` — tests existentes + nuevos casos con `--field`
