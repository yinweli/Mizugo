# CLAUDE.md

This file provides guidance for AI CLIs (such as Claude Code or ChatGPT Codex) when working in this repository.

## Language Conventions

- Documentation and comments: Traditional Chinese
- Code naming (variables, functions, classes): English
- Implementation plans: Chinese descriptions with English technical terms

## Workflow

- Use Python scripts instead of long shell command chains for multi-step, repetitive, or cross-file tasks.
- Use Python first for:
  - batch file operations
  - structured data parsing or conversion
  - code generation
  - any task involving loops, conditions, or cross-file text processing

- Generate mock data in batches according to project conventions instead of writing test data manually one by one.
- After running tests or lint checks, summarize only failures and warnings from the raw output.

## Context Management

- Manage context usage proactively during long conversations or when working with large files.
- In long conversations, periodically summarize key decisions, current status, and remaining tasks.
- When working with files longer than 500 lines, handle them in focused sections instead of loading everything at once.
- Before complex multi-step tasks, briefly restate the relevant code style rules and project constraints.
- If context becomes overloaded, warn the user and suggest starting a new conversation or consolidating the current state into a file.

## Code Style

- In Go and C#, `false` checks must be explicit. Negation-style checks are forbidden.
  - Example: `if x == false {}` / `if (x == false) {}`
- In TypeScript, negation-style checks such as `if (!x) {}` are allowed.

- Closing comments are only allowed for control flow blocks: `if`, `for`, `switch`
  - Example: `} // if`
- Closing comments on functions or methods are forbidden.

- Variable names must always be singular, even for slices, arrays, maps, and other collections.
  - Example: `item := []Item{}`, `hero := []Hero{}`
  - Forbidden: `items := []Item{}`, `heroes := []Hero{}`

- Iterator naming:
  - Use `itor` for general iteration
  - Use `k, v` for map iteration

## Commit Message Convention

Format: `Type | Description`

| Type      | Usage                          |
| :-------- | :----------------------------- |
| `Feature` | New feature or update          |
| `Fix`     | Bug fix                        |
| `Sheet`   | Update sheet data              |
| `Message` | Update proto messages          |
| `Doc`     | Update documentation           |

- Description must be in Traditional Chinese.
- Do not invent new type names (e.g. `Bugfix`, `Update`, `Refactor` are all forbidden).

## Project Overview

Mizugo is a game server framework written in Go (module: `github.com/yinweli/Mizugo/v2`). It provides TCP networking, message processing, a two-tier database layer (Redis + MongoDB), an entity/module system, and configuration management. Documentation and comments are in Traditional Chinese.

## Development / Build / Common Commands

All commands use [Task](https://taskfile.dev/) (install per official docs).

```bash
task lint                # Format and lint (golangci-lint fmt + run, markdownlint, prettier)
task test                # Run tests (not defined in Taskfile; use go test directly)
task proto               # Regenerate protobuf files
task db                  # Start Redis + MongoDB via Docker
task install             # Install all dev tools (golangci-lint, buf, csharpier, etc.)
```

### Running Tests Directly

```bash
# All tests (excluding support/ directory)
go test $(go list ./... | grep -v "support")

# Single package
go test ./mizugos/entitys/...

# Single test
go test ./mizugos/entitys/... -run TestEntityName

# With coverage
go test -coverprofile=coverage.txt -covermode=atomic $(go list ./... | grep -v "support")
```

`redmos`, `trials`, and integration tests require Redis and MongoDB to be running (use `task db` to start them).

### Lint

Uses golangci-lint v2 with 47+ linters enabled. `goimports` local import prefix: `github.com/yinweli/Mizugo/v2`.

## Architecture

### Manager-Based Design

The framework is accessed through six singleton managers in the `mizugos` package:

| Manager  | Package             | Accessor         | Purpose                              |
| -------- | ------------------- | ---------------- | ------------------------------------ |
| Config   | `configs.Configmgr` | `mizugos.Config` | Viper-based configuration management |
| Logger   | `loggers.Logmgr`    | `mizugos.Log`    | Logging (Zap-based)                  |
| Network  | `nets.Netmgr`       | `mizugos.Net`    | TCP networking and session management|
| Database | `redmos.Redmomgr`   | `mizugos.Redmo`  | Redis (cache) + MongoDB (persistence)|
| Entity   | `entitys.Entitymgr` | `mizugos.Entity` | Game objects with module/event system|
| Pool     | `pools.Poolmgr`     | `mizugos.Pool`   | Goroutine pool (ants)                |

Server lifecycle: `mizugos.Start()` → run main loop → `mizugos.Stop()`.

### Package Layering

Each layer follows strict dependency rules: lower layers may reference upper layers, but **never the reverse**. Packages within the same layer must not reference each other.

```text
Test layer:      testdata, mizugos/trials
Utility layer:   mizugos/ctxs, mizugos/helps, mizugos/msgs
Common layer:    mizugos/cryptos, mizugos/iaps, mizugos/nets, mizugos/pools, mizugos/procs
Component layer: mizugos/configs, mizugos/entitys, mizugos/loggers, mizugos/redmos
```

### Message Processing

The `procs` package provides three processor types: JSON, Proto (protobuf), and Raven (custom binary format). Each processor handles encoding/decoding through a configurable codec chain. Messages must include an `int32 messageID` field.

### Database (redmos)

Uses a two-tier architecture: Redis as the cache layer ("Major") and MongoDB as the persistence layer ("Minor"), with "Mixed" operations spanning both. The package provides 60+ command types supporting get/set, queues, locks, aggregation, and more.

### Entity System

Entities support modular composition via the `Module` interface (with Awake/Start lifecycle), an event system (one-shot, delayed, periodic), and message handler registration. Sessions bind network connections to entities.

## Key Paths

- Framework source: `mizugos/`
- Protobuf definitions: `support/proto-mizugo/`, `support/proto-test/`
- Test server example: `support/test-server/`
- Unity client: `support/client-unity/`
- Test data: `testdata/`

## Go Version

Requires Go 1.25.0+ and Protocol Buffers v3.
