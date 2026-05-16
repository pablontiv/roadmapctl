# Auto-update — Staged Async Pattern

roadmapctl updates itself automatically using a **staged async** pattern: the new binary is downloaded in the background during run N, and applied at the start of run N+1. The update is transparent — no prompts, no interruptions, no downtime.

## Flow per Invocation

```
Run N:
  1. ApplyStagedIfAvailable()  — sync: checks ~/.cache/roadmapctl/staged/ for a newer binary
     → if found: atomic rename over current binary, re-exec (process is replaced)
     → if not found or not newer: continues normally
  2. go FetchAndStage(version) — goroutine: downloads next release in background
     → writes to ~/.cache/roadmapctl/staged/<tag>/roadmapctl
     → exits silently on any error

Run N+1:
  → staged binary is detected, applied, and roadmapctl re-execs with the new version
```

The user-visible effect is that the version shown by `roadmapctl --version` changes on the second invocation after a new release.

## Staging Directory

```
~/.cache/roadmapctl/staged/
  v0.2.0/
    roadmapctl          # Linux / macOS
    roadmapctl.exe      # Windows
```

The directory is created by the downloader. Each staged release occupies its own version subdirectory. Staging directories are cleaned up automatically when the binary is applied.

## Behavior by OS

| Concern | Unix (Linux / macOS) | Windows |
|---------|---------------------|---------|
| Apply strategy | `os.Rename(staged, current)` — atomic on same filesystem | rename current → `.old`, copy staged → current, remove `.old` |
| Re-exec | `syscall.Exec` — replaces process in same PID, transparent to caller | `exec.Command` + `os.Exit(0)` — new process launched, old exits |
| Binary in use | Not an issue — rename is atomic | Rename of in-use binary fails; copy avoids the lock |

## Escape Hatches

| Method | Effect |
|--------|--------|
| `ROADMAPCTL_NO_UPDATE=1` | Skips both `ApplyStagedIfAvailable` and `FetchAndStage` entirely |
| `version == "dev"` | Auto-update is always disabled for development builds |

## Silent Failure Policy

All auto-update errors are suppressed silently — the current command is never interrupted:

- **Network errors** during download: silently skip; staging directory is not created.
- **SHA256 mismatch**: returns an error internally; no file is written to the staging directory.
- **Permission errors** applying the update (e.g., `/usr/local/bin` owned by root): silently skip; the command proceeds normally.

If the update is consistently failing silently, check:

```bash
ls ~/.cache/roadmapctl/staged/     # see what's staged
roadmapctl --version               # current version
ROADMAPCTL_NO_UPDATE=1 roadmapctl --version  # verify the binary works without update logic
```
