# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Mizugo is a game server framework in Go (module: `github.com/yinweli/Mizugo/v2`). It provides TCP networking, message processing protocols, dual-layer database abstraction (Redis/MongoDB), entity-component systems, and utility packages. Code comments and documentation are in Traditional Chinese.

## Commands

### Task runner (preferred)

```bash
task lint    # Format code (gofmt, goimports) and run golangci-lint
task test    # Run all unit tests with coverage
task bench   # Run benchmark tests
task db      # Start Redis and MongoDB via Docker (required for integration tests)
task stop    # Stop Docker containers
task proto   # Generate protobuf code
task install # Install dev tools (golangci-lint, buf, goimports, etc.)
```

### Direct Go commands

```bash
# Run all tests
go test ./... -cover

# Run tests for a specific package
go test ./mizugos/configs -v

# Run a single test
go test ./mizugos/configs -run TestConfigmgr

# Run with race detection
go test -race ./mizugos/configs

# Run linter manually
golangci-lint run --color always
```

### Integration test requirement

Tests in `redmos` and `trials` require live Redis 6+ and MongoDB 6.0+ instances. Use `task db` to start them before running those tests.

## Architecture

### Dependency layer rules

Packages are organized into strict layers. **Lower layers must never import upper layers. Same-layer packages must not import each other.**

```
Layer 4 – Components:  configs, entitys, loggers, redmos
Layer 3 – General:     cryptos, iaps, nets, pools, procs
Layer 2 – Tools:       ctxs, helps, msgs
Layer 1 – Test:        testdata, trials
```

All core packages live under `mizugos/`. The `support/` directory contains standalone programs (test servers, test clients, proto definitions) that are not part of the framework itself.

### Framework entry point (`mizugos/mizugo.go`)

`mizugos.Start()` / `mizugos.Stop()` initialise and tear down all managers. Global manager singletons:

- `mizugos.Config` – Configmgr (Viper-based, supports files/strings/env)
- `mizugos.Logger` – Logmgr (named loggers, ZapLogger or EmptyLogger)
- `mizugos.Network` – Netmgr (TCP listener/connecter)
- `mizugos.Redmo` – Redmomgr (Redis + MongoDB dual-layer)
- `mizugos.Entity` – Entitymgr (entity/module/event system)
- `mizugos.Pool` – Poolmgr (goroutine pool via `ants`)

### Networking (`nets`)

- **Listener** – accepts TCP connections; **Connecter** – initiates TCP connections.
- **Sessioner** – per-connection session with lifecycle: `Start → EventStart → recv/send loops → EventStop → Finalize`.
- **Codec chain** – encoding pipeline; encode runs forward, decode runs reverse. Each codec can be stacked.
- Default limits: HeaderSize = 4 bytes, PacketSize = 65535 bytes, ChannelSize = 1000.
- Session callbacks: `Bind` (init), `Unbind` (cleanup), `Publish` (event dispatch), `Wrong` (error handler).

### Message processors (`procs`)

Implement the `Processor` interface (`Encode`/`Decode`/`Process`/`Add`/`Del`/`Get`). Built-in processors:

- **Json** – JSON serialization
- **Proto** – Protocol Buffers serialization
- **Raven** – custom protocol serialization

### Entity system (`entitys`)

- **Entity** holds an ID, a set of **Modules**, and **Events**.
- Module lifecycle: `Awaker.Awake()` then `Starter.Start()`.
- Built-in events: `EventDispose`, `EventShutdown`, `EventRecv`, `EventSend`.
- Events can be single-shot, delayed, or recurring.

### Database (`redmos`)

- **Major** – Redis layer for caching and fast queries.
- **Minor** – MongoDB layer for persistence and complex queries.
- **Mixed** – combines both layers. Operations are composed via `Submittor` interface.

### Error handling (`helps`)

Use `helps.Err` for structured errors with numeric codes. Format: `<FunctionName>: <message> (<error_code>)`. Predefined codes: `Success (0)`, `ErrUnknown (1)`, `ErrUnwrap (2)`.

### Test utilities (`trials`)

```go
trials.SetupRedis(uri)           // connect to test Redis
trials.SetupMongo(uri, dbName)   // connect to test MongoDB
trials.ProtoBuild(...)           // compile proto files during tests
```

## Key dependencies

| Library | Purpose |
|---|---|
| `go.uber.org/zap` | Logging |
| `github.com/spf13/viper` | Configuration |
| `github.com/panjf2000/ants/v2` | Goroutine pool |
| `github.com/redis/go-redis/v9` | Redis client |
| `go.mongodb.org/mongo-driver` | MongoDB client |
| `google.golang.org/protobuf` | Protocol Buffers |
| `github.com/stretchr/testify` | Test assertions |

## Linter configuration

`.golangci.yml` enables 25+ linters. Notable limits: line length 200, function length 200 lines / 150 statements, cyclomatic complexity 50, duplication threshold 400 LOC. Test files are exempt from MND and dupl checks.

## Code Style Requirements

### Boolean Checks - **Strict Rule**

**MUST** use explicit `false` checks, while `true` checks use standard form.

```go
if x == false { ... }  // Required: explicit false check
if x { ... }           // Standard: check for true
// Forbidden: if !x
```

### Block Ending Comments - **Strict Rule**

Reserve ending comments for control flow only (`if`, `for`, `switch`). **NEVER** add ending comments for functions/methods.

```go
if condition {
    // logic
} // if

for i := range items {
    // logic
} // for

switch x {
case 1:
    // logic
} // switch

// NEVER add ending comments for functions/methods
```

### Iterator Naming

```go
for itor := range items {
    // logic
} // for

for k, v := range someMap {
    // logic
} // for
```

- `itor`: Default iterator variable name
- `k, v`: Use when iterating over Maps

### Language Conventions

- **Documentation & Comments**: Traditional Chinese
- **Code Naming (variables, functions)**: English
- **Implementation Plans**: Chinese description + English technical terms

## Interaction Guidelines

Communicate with users in Traditional Chinese
