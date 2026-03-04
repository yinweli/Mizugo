# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common guildlines

@CLAUDE-common guildlines.md

## Project Overview

Mizugo is a game server framework written in Go (module: `github.com/yinweli/Mizugo/v2`). It provides TCP networking, message processing, dual-layer database (Redis + MongoDB), entity/module systems, and configuration management. Documentation and comments are in Traditional Chinese.

## Common Commands

All commands use [Task](https://taskfile.dev/) (install via its official docs).

```bash
task lint                # Format and lint (golangci-lint fmt + run, markdownlint, prettier)
task test                # Run tests (not defined in Taskfile; use go test directly)
task bench               # Performance benchmarks
task proto               # Regenerate protobuf files
task db                  # Start Redis + MongoDB via Docker
task install             # Install all dev tools (golangci-lint, buf, csharpier, etc.)
```

### Running tests directly

```bash
# All tests (excludes support/ directory)
go test $(go list ./... | grep -v "support")

# Single package
go test ./mizugos/entitys/...

# Single test
go test ./mizugos/entitys/... -run TestEntityName

# With coverage
go test -coverprofile=coverage.txt -covermode=atomic $(go list ./... | grep -v "support")
```

Tests for `redmos`, `trials`, and integration tests require running Redis and MongoDB (use `task db` to start them).

### Linting

Uses golangci-lint v2 with 47+ linters. Local import prefix for goimports: `github.com/yinweli/Mizugo/v2`.

## Architecture

### Manager-based Design

The framework is accessed through six singleton managers via the `mizugos` package:

| Manager | Package | Access | Purpose |
| ------- | ------- | ------ | ------- |
| Config | `configs.Configmgr` | `mizugos.Config` | Viper-based configuration |
| Logger | `loggers.Logmgr` | `mizugos.Log` | Logging (Zap-based) |
| Network | `nets.Netmgr` | `mizugos.Net` | TCP networking & sessions |
| Database | `redmos.Redmomgr` | `mizugos.Redmo` | Redis (cache) + MongoDB (persistence) |
| Entity | `entitys.Entitymgr` | `mizugos.Entity` | Game objects with module/event systems |
| Pool | `pools.Poolmgr` | `mizugos.Pool` | Goroutine pool (ants) |

Server lifecycle: `mizugos.Start()` → run loop → `mizugos.Stop()`.

### Package Layer Hierarchy

Layers enforce a strict dependency rule: lower layers may reference upper layers, but **not** vice versa. Same-layer packages cannot reference each other.

```text
Test Layer:     testdata, mizugos/trials
Tool Layer:     mizugos/ctxs, mizugos/helps, mizugos/msgs
Common Layer:   mizugos/cryptos, mizugos/iaps, mizugos/nets, mizugos/pools, mizugos/procs
Component Layer: mizugos/configs, mizugos/entitys, mizugos/loggers, mizugos/redmos
```

### Message Processing

Three processor types (`procs` package): JSON, Proto (protobuf), and Raven (custom binary). Each handles encode/decode with a configurable codec chain. Messages require an `int32 messageID` field.

### Database (redmos)

Dual-layer architecture: Redis as cache ("Major"), MongoDB as persistence ("Minor"), with "Mixed" operations spanning both. The package has 60+ command types for get/set, queues, locks, aggregation, etc.

### Entity System

Entities support modular composition via `Module` interface (Awake/Start lifecycle), an event system (single-shot, delayed, periodic), and message handler registration. Sessions bind network connections to entities.

## Key Paths

- Framework source: `mizugos/`
- Protobuf definitions: `support/proto-mizugo/`, `support/proto-test/`
- Test server example: `support/test-server/`
- Unity client: `support/client-unity/`
- Test data: `testdata/`

## Go Version

Requires Go 1.25.0+ and Protocol Buffers v3.
