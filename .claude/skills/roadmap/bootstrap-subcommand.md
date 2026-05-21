# /roadmap bootstrap

Fuerza re-escaneo del repo y wizard de adopción de gitflow, even if `.roadmapctl.toml` already has values.

## Flujo

1. Ejecutar `roadmapctl bootstrap --repo <repo> --output json` para detectar el estado actual.
2. Iniciar el wizard de adopción gitflow (ver [bootstrap-reference.md](bootstrap-reference.md) sección "Wizard de adopción gitflow").
3. El escaneo es abierto — el LLM decide qué leer según lo que encuentra en el repo.
4. Tras confirmación del usuario, escribir los style fields bajo `[gitflow]` en `.roadmapctl.toml`.
5. Re-ejecutar `roadmapctl bootstrap` para confirmar los valores escritos.
6. **Terminar**: no ejecutar plan/loop/pending/decision tras el bootstrap subcommand.

## Nota

Este subcomando solo gestiona la configuración gitflow. No inicia ejecución de tasks ni materialización de roadmap.
