---
estado: Pending
tipo: task
---
# T003: Configure squash-only merge and delete branch on merge

**Contribuye a**: align merge strategy with the ecosystem standard (squash-only keeps a clean linear history; auto-delete avoids stale branch accumulation).

## Alcance

**In**:
- PATCH repos/pablontiv/roadmapctl: `squash_merge: true`, `merge_commit: false`, `rebase_merge: false`, `delete_branch_on_merge: true`

**Out**:
- No file changes

## Criterios de Aceptación

- `gh repo view pablontiv/roadmapctl --json squashMergeAllowed,mergeCommitAllowed,rebaseMergeAllowed,deleteBranchOnMerge` returns `{"squashMergeAllowed":true,"mergeCommitAllowed":false,"rebaseMergeAllowed":false,"deleteBranchOnMerge":true}`

## Fuente de verdad

- GitHub API: PATCH repos/pablontiv/roadmapctl
