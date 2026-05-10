# Decision: keep `roadmapctl materialize` as deterministic writer

Status: Accepted
Date: 2026-05-10

## Decision

Keep `roadmapctl materialize` as the required deterministic writer for `/roadmap plan` materialization. The roadmap skill should continue to create a structured temp plan, run `roadmapctl materialize --dry-run`, show only concise dry-run evidence by default, and apply through roadmapctl-owned batch apply after approval.

Do not reintroduce prompt-side direct roadmap file writes as the normal flow. Do not remove `roadmapctl materialize` in favor of a guard-only command yet.

## Evidence considered

- Local docs already define materialize as a deterministic writer that rejects free-form prose and owns numbering, dependency links, canonical paths, dry-run, frozen change-set apply, and postcheck recovery (`docs/materialize-plan-schema.md`, `docs/roadmap-skill-integration.md`).
- Backscroll evidence shows the original consistency risk: invalid single-summary fallback files were explicitly tested and blocked with `RMC_STRUCTURE_SINGLE_FILE_FALLBACK` in prior roadmapctl sessions.
- Backscroll evidence also shows current UX cost: materialize dry-runs previously created excessive context from `changes[].content`/diffs, which was addressed by the token-light dry-run flow.
- Headless regression coverage now verifies no-argument `/roadmap` routes through `roadmapctl pending` and materialize dry-run can be reviewed concisely without repository writes.

## Alternatives

### Keep deterministic writer

Pros:
- Preserves the strongest guard against `*-tasks.md` fallback and malformed roadmap structure.
- Keeps numbering, canonical paths, dependency serialization, preconditions, and postcheck in code instead of prompt logic.
- Supports dry-run review, frozen change-set apply, and recovery after partial apply/postcheck failure.

Cons:
- More implementation and documentation surface area than a simple validator.
- Requires token discipline around dry-run output.

### Guard-only validation

Pros:
- Smaller CLI write surface.
- Easier to reason about if the skill owns all writes.

Cons:
- Reopens prompt-side writer drift: the skill would need to duplicate numbering, path generation, dependency links, and file rendering.
- Validation would catch some failures only after files were already written.
- Historical single-file fallback failures would be easier to reintroduce.

### Direct write plus `roadmapctl check`

Pros:
- Lowest CLI complexity.
- Fastest for simple plans.

Cons:
- Weakest safety model. It relies on the agent to generate every file correctly and then recover after failures.
- Increases chance of non-canonical intermediate states and token-heavy manual diffs.
- Conflicts with current skill guidance that deterministic writes belong to roadmapctl.

## Revisit criteria

Reconsider simplifying or deprecating `roadmapctl materialize` only if all of the following are true:

1. Headless and fixture coverage prove direct or guard-only flows cannot create single-summary fallback files, stale schema output, invalid dependency links, or non-canonical paths.
2. The skill no longer needs to duplicate deterministic writer logic in prompt text.
3. Token overhead remains low without hiding required diagnostics or recovery evidence.
4. Postcheck failure recovery is at least as explicit as the current frozen change-set flow.
5. A migration plan exists for existing docs, tests, and operator muscle memory.

Until then, optimize materialize UX and schema discoverability rather than replacing the writer.
