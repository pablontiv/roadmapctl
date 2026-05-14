# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Releases are automated via CI from conventional commits.

## [Unreleased]

### Added

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
