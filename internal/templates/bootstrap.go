package templates

const BaseStemContent = `version: 2
scope:
  match: "*.md"

schema:
  estado:
    type: enum
    required:
      match: ["T*"]
    match: ["O*", "T*"]
    values: [Pending, Specified, In Progress, Completed, Blocked, On Hold, Obsolete]

  tipo:
    type: enum
    required:
      match: ["O*", "T*"]
    match: ["O*", "T*"]
    values: [outcome, task]

  id:
    type: sequence
    match:
      "O*": { prefix: O, digits: 2 }
      "T*": { prefix: T, digits: 3 }

links:
  blocked_by:
    target: '^(\./|\.\./|.*/)T[0-9]{3}-[^/]+\.md$'
  reference:
    target: ".*"

validate:
  - field: tipo
    rule: non_empty
`

const DefaultRoadmapctlTOML = `done_statuses = ["Completed", "Obsolete"]
active_statuses = ["Pending", "Specified", "In Progress"]
leaf_filter = "isIndex == false"
outcome_close_verify = []
pr_merge_strategy = "squash"
commit_style = "conventional"
auto_push = true
required_code_coverage = 85.0
loop_max_tasks = 0
parallel = true
autonomy = "until_done"
compact_after_task_commit = true
pr_mode = false

[status_values]
pending = "Pending"
specified = "Specified"
in_progress = "In Progress"
completed = "Completed"
blocked = "Blocked"
obsolete = "Obsolete"
`
