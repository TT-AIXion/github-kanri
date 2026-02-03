# Product Requirements (Summary)

## Background

- Local GitHub repo management is scattered and manual
- Agent-driven runs need safety and reproducibility
- Shared assets/skills sync is painful without automation

## Goals

- Safe, non-interactive CLI operations
- Fast bulk management of local repos and skills
- Deterministic behavior with clear failure reasons

## Target Users

- Primary: developers, AI agents
- Secondary: operators, PMs

## Must Have

- macOS local CLI
- Single JSON config
- Allow/deny safety guard
- Repo scan/status/bulk exec
- Skills sync for shared assets

## Should Have

- Fail on conflicts by default, `--force` to override
- Ambiguous matches return candidates only
- Glob include/exclude

## Could Have

- Watch/schedule
- Verify/diff

## Non-goals

- GUI or web UI
- Natural language input
- External services (GitHub API)
- Persistent logging

## Success Metrics

- Repo/skills ops complete in seconds
- Zero accidental destructive actions
- 100% reproducibility in LLM runs

## Constraints

- Local-only operation
- Default under `reposRoot`
- Single config file
