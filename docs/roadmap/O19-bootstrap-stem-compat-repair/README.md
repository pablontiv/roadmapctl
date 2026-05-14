---
tipo: outcome
---
# Bootstrap stem compat repair

`roadmapctl bootstrap` detecta cuando el `.stem` de un repo tiene un schema incompatible
(estado requerido en Outcomes) y ofrece al usuario repararlo interactivamente, en lugar de
solo bloquear y reportar.

El resultado observable es que un agente ejecutando `/roadmap plan` en un repo con `.stem`
legacy puede resolver el bloqueo desde el mismo flujo de bootstrap, sin intervención manual
ni copia de archivos entre repos.
