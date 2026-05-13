---
tipo: outcome
---
# O16: CI verde — segunda ola

Los tasks de O14/O15 están Completed pero la pipeline sigue sin producir releases. La investigación reveló 8 bugs distintos en 5 jobs fallidos. Este outcome cierra los bugs restantes para que `auto-tag` de crossbeam pueda correr por primera vez.

Jobs que deben quedar verdes: `ci/Lint`, `ci/Test`, `smoke/ubuntu-latest`, `smoke/macos-latest`, `smoke/windows-latest`.
