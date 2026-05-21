---
tipo: outcome
---
# O30: Materialización paralela en `/roadmap plan`

Reemplazar el modelo singleton-writer del skill `/roadmap plan` por un coordinador que despacha múltiples Agents en paralelo tras aprobación del usuario. Cada Agent escribe un subset de archivos (1 Write por archivo), reduciendo la latencia de materialización y habilitando el patrón de waves (Outcomes primero, Tasks en paralelo).
