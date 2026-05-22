# /roadmap bootstrap

Fuerza re-escaneo del repo y wizard de adopción de gitflow, even if `.roadmapctl.toml` already has values.

## Flujo

1. Ejecutar `roadmapctl bootstrap --repo <repo> --output json` para detectar el estado actual.
2. Si `missing_settings` o `empty_settings` contienen campos de `gitflow`, ejecutar el wizard de adopción gitflow (ver [bootstrap-reference.md](bootstrap-reference.md) sección "Wizard de adopción gitflow").
3. El escaneo es abierto — el LLM decide qué leer según lo que encuentra en el repo.
4. Presentar dry-run del bloque `[gitflow]` propuesto. Solo escribir tras confirmación explícita del usuario.
5. Tras confirmación, escribir los campos bajo `[gitflow]` en `.roadmapctl.toml`.
6. Re-ejecutar `roadmapctl bootstrap --repo <repo> --output json` para confirmar los valores escritos.
7. **Terminar**: no ejecutar plan/loop/pending/decision tras el bootstrap subcommand.

## Nota

Este subcomando gestiona la configuración gitflow. `roadmapctl bootstrap` es ahora read-only;
las escrituras (crear archivos, reparar `.stem`) se hacen con `roadmapctl bootstrap init --apply`.
