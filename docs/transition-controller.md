# Transition controller contract

`roadmapctl transition` is the planned policy layer for task status changes. It decides whether a roadmap task can move between operational status roles, then emits the exact Rootline change that would be applied by a later `--apply` flow.

The controller is deterministic and read-model driven. Concrete status strings always come from `.roadmapctl.toml`/legacy config plus the effective Rootline schema; rules must not hardcode labels such as `Pending` or `Completed` except as built-in config defaults.

## Actions

| Action | Mutates | Meaning |
|--------|---------|---------|
| `can-start` | no | Report whether a task can move to the configured `in_progress` role. |
| `can-complete` | no | Report whether a task can move to the configured `completed` role. |
| `start` | dry-run by default | Plan `estado=<status_values.in_progress>` for an allowed start. |
| `complete` | dry-run by default | Plan `estado=<status_values.completed>` for an allowed completion. |
| `set-status` | dry-run by default | Plan an explicit `estado=<value>` when the value is valid in schema and policy permits it. |

Future apply flows must remain explicit (`--apply`) and run a postcheck after Rootline mutation.

## Required operational roles

The config must provide these role values:

- `status_values.pending`
- `status_values.specified`
- `status_values.in_progress`
- `status_values.completed`
- `status_values.blocked`
- `status_values.obsolete`

`done_statuses` determines dependency satisfaction. `active_statuses` determines candidates listed by read-only commands. The transition controller must use these lists instead of raw status literals.

## Core rules

### `can-start`

Allowed when:

1. the target path resolves to a task in the normalized read model;
2. current `estado` is in `active_statuses`;
3. current `estado` is not in `done_statuses`;
4. all `blocked_by` dependencies resolve to tasks whose `estado` is in `done_statuses`;
5. the concrete target status `status_values.in_progress` exists in the effective Rootline schema.

Blocked when any dependency is missing or not done. The response includes each blocker path and observed status.

### `can-complete`

Allowed when:

1. the target path resolves to a task;
2. current `estado` is not in `done_statuses`;
3. the concrete target status `status_values.completed` exists in schema.

`can-complete` does not execute acceptance checks, tests, commits, or PR work. The caller remains responsible for task-specific verification before applying completion.

### `set-status`

Allowed when:

1. target path resolves to a task;
2. requested status exists in schema;
3. setting the status would not bypass dependency policy for the `in_progress` role.

Setting to a done role is allowed by policy after caller verification. Setting to a non-role valid schema value such as `On Hold` is allowed, reported as `role: "custom"`, and excluded from active/done semantics unless config maps it into a role.

## `On Hold` and other valid non-role statuses

`On Hold` may exist in `.stem` without appearing in operational roles. In that case:

- it is schema-valid;
- it is not dependency-satisfying unless included in `done_statuses`;
- it is not listed as active unless included in `active_statuses`;
- `can-start` is blocked because the task is outside active roles;
- `set-status --status "On Hold"` is allowed as an explicit custom status.

This same behavior applies to any project-specific valid status that is not mapped to an operational role.

## JSON shape

```json
{
  "version": 1,
  "kind": "roadmapctl/transition",
  "summary": { "status": "ok", "errors": 0, "warnings": 0, "infos": 0 },
  "root": "/repo",
  "roadmap_root": "/repo/docs/roadmap",
  "action": "can-start",
  "path": "O01-work/T001-task.md",
  "allowed": true,
  "current_status": "Pending",
  "target_status": "In Progress",
  "role": "in_progress",
  "reasons": ["all dependencies are done"],
  "blockers": [],
  "changes": [
    { "path": "O01-work/T001-task.md", "field": "estado", "from": "Pending", "to": "In Progress" }
  ],
  "diagnostics": []
}
```

When blocked, `allowed` is `false`, `changes` is empty, and `reasons`/`blockers` explain the deterministic cause.

## Diagnostics

| ID | Severity | Meaning |
|----|----------|---------|
| `RMC_TRANSITION_TASK_NOT_FOUND` | error | Requested path is not a known task in the read model. |
| `RMC_TRANSITION_STATUS_UNKNOWN` | error | Requested or configured target status is absent from the effective schema. |
| `RMC_TRANSITION_DEPENDENCY_BLOCKED` | warning | Start is blocked by at least one dependency outside `done_statuses`. |
| `RMC_TRANSITION_ROLE_MISSING` | error | Required operational role is missing from config or empty. |
| `RMC_TRANSITION_NOT_ACTIVE` | warning | Task is outside `active_statuses` for a start action. |
| `RMC_TRANSITION_ALREADY_DONE` | warning | Task is already in a `done_statuses` role for start/complete. |
| `RMC_TRANSITION_APPLY_FAILED` | error | Future apply mode failed during Rootline mutation or postcheck. |

Warnings make dry-run/can-* output informative without implying invalid roadmap state. `--strict` promotes warnings to validation exit code `1`, matching global diagnostics policy.
