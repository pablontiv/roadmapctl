---
estado: Pending
tipo: task
---
# T002: Configure branch protection on master

**Contribuye a**: enforce the ecosystem security baseline — no commit reaches master without a PR review, mirroring crossbeam and rootline.

## Alcance

**In**:
- Configure branch protection on `master` via GitHub API:
  - `required_pull_request_reviews.required_approving_review_count: 1`
  - `required_pull_request_reviews.dismiss_stale_reviews: true`
  - `allow_force_pushes: false`
  - `allow_deletions: false`

**Out**:
- No file changes

## Criterios de Aceptación

- `gh api repos/pablontiv/roadmapctl/branches/master/protection` returns HTTP 200
- `required_pull_request_reviews.required_approving_review_count` >= 1
- `allow_force_pushes.enabled` = false
- `allow_deletions.enabled` = false

## Fuente de verdad

- GitHub API: PUT repos/pablontiv/roadmapctl/branches/master/protection
