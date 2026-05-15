---
tipo: outcome
---
# O21: Auto-update staged async

roadmapctl detecta automáticamente nuevas versiones en GitHub Releases y se actualiza sin bloquear el comando actual. En la siguiente invocación, el binario actualizado entra en uso con re-exec transparente.

Cada invocación lanza en background la descarga de la nueva versión al directorio de staging (`~/.cache/roadmapctl/staged/<version>/`). Al arrancar, si hay un binario staged más nuevo, se aplica con atomic rename y re-exec antes de ejecutar el comando pedido.
