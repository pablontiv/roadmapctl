# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Releases are automated via CI from conventional commits.

## [Unreleased]

### Added (O28)

- `[gitflow]` config section with three style fields: `branch_style`, `pr_title_style`, `pr_body_style`
- `roadmapctl bootstrap` now serializes gitflow fields in JSON output
- New diagnostic `RMC_GITFLOW_NOT_CONFIGURED` (info) when gitflow is unconfigured
- New skill subcommand `/roadmap bootstrap` for on-demand gitflow setup wizard

### Changed (O28)

- `/integrate` skill: Fase 2 always runs (was conditional on `pr_mode`); branch/commit/PR generation uses TOML style fields
- `/roadmap loop` paso 9 declares `/integrate` as the only gate to git/gh

### Deprecated (O28)

- `pr_merge_strategy` TOML field: now emits `RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED` warning; use `[gitflow]` fields

### Added (Previous)

- `LICENSE`, `CONTRIBUTING.md`, `CHANGELOG.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md` (ecosystem documentation baseline)
- `CODEOWNERS` (`* @pablontiv`) to enforce review requirements
- `dependabot.yml` for automated dependency updates (gomod + github-actions)
- CodeQL and OpenSSF Scorecard workflows via crossbeam reusable workflows
- Squash-only merge + delete-branch-on-merge enabled
- Repo made public; branch protection on `master` requiring PR review

## [v0.0.1] - 2026-04

### Added

- Initial release: `doctor`, `check`, `lint` guards
- Read-only state: `context`, `pending`, `next`, `decision`
- Controlled mutation: `transition`, `materialize`, `bootstrap`
- `/roadmap` skill distributed via pre-push hook
- Stable diagnostics and exit codes per CLI contract
