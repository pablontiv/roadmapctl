---
estado: Completed
tipo: task
---
# T004: Cargar configuración de roadmap y resolver roots

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE1 y CE2

[[blocked_by:./T002-create-go-cli-skeleton.md]]

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: parsear config localmente o vía Rootline sin mover lógica al motor.
- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: config inválida produce diagnostics, no escritura.

## Contexto

El skill `/roadmap` usa `.claude/roadmap.local.md` para `roadmap-root`, estados, filtros y opciones operacionales. `roadmapctl` debe resolver esa configuración de forma determinística y cross-platform.

## Alcance

**In**:
1. Leer `.claude/roadmap.local.md` desde `--repo` o cwd.
2. Parsear frontmatter YAML.
3. Aplicar defaults documentados cuando correspondan.
4. Resolver `roadmap-root` como path absoluto contenido dentro del repo.
5. Detectar escape por `..` o symlinks según política definida.
6. Soportar override explícito `--roadmap-root`.

**Out**:
- Workspace mode completo si no entra en MVP.
- Validación profunda de schema Rootline.
- Creación de config faltante.

## Estado inicial esperado

- Skeleton CLI existe.
- Diagnostics model existe o está en progreso.

## Criterios de Aceptación

- Fixture válido resuelve `roadmap-root` correctamente.
- Fixture con `roadmap-root: ../outside` falla con diagnostic de path escape.
- Config faltante produce error claro y exit `2` o diagnostic configurado.
- Tests cubren separadores Windows y paths relativos.

## Fuente de verdad

- `.claude/roadmap.local.md`
- `internal/config/`
- `internal/fsx/`
- `docs/roadmap/O01-roadmapctl-mvp/README.md`
