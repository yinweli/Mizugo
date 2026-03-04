# Language Conventions

| Aspect | Convention |
| ------ | ---------- |
| Documentation & Comments | Traditional Chinese |
| Code Naming (variables, functions, classes) | English |
| Implementation Plans | Chinese description + English technical terms |

# Code Style Requirements

## 1. Boolean Checks — Strict Rule

**MUST** use explicit `false` checks; `true` checks use standard form.

```go
if x == false { ... }  // Required: explicit false check
if x { ... }           // Standard: check for true
// Forbidden: if !x
```

```csharp
if (x == false) { ... }  // Required: explicit false check
if (x) { ... }           // Standard: check for true
// Forbidden: if (!x)
```

```typescript
if (x === false) { ... }  // Required: explicit false check
if (x) { ... }            // Standard: check for true
// Forbidden: if (!x)
```

## 2. Block Ending Comments — Strict Rule

Reserve ending comments for control flow blocks only (`if`, `for`, `switch`).
**NEVER** add ending comments for functions/methods.

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
func DoSomething() {
    // logic
}
```

```csharp
if (condition)
{
    // logic
} // if

for (int i = 0; i < count; i++)
{
    // logic
} // for

switch (x)
{
    case 1:
        // logic
        break;
} // switch
```

```typescript
if (condition) {
    // logic
} // if

for (const itor of items) {
    // logic
} // for

switch (x) {
    case 1:
        // logic
} // switch
```

## 3. Variable Naming — Strict Rule

**MUST NOT** use plural forms for variable names. Use singular forms consistently.

```go
// Correct: singular form
item := []Item{...}
hero := []Hero{...}
data := []Data{...}

// Forbidden: plural forms
items := []Item{...}
heroes := []Hero{...}
datas := []Data{...}
```

## 4. Iterator Naming

```go
// Default iterator: use "itor"
for itor := range item {
    // logic
} // for

// Map iteration: use "k, v"
for k, v := range someMap {
    // logic
} // for
```

# AI Workflow Optimization

When executing shell commands (e.g., batch file operations, text transformations, data processing, code generation), prefer using **Python scripts** over long chains of shell commands or repeated tool invocations. Python scripts can accomplish complex multi-step tasks in a single execution, significantly reducing token consumption and round-trip overhead.

**Example — batch renaming files:**

```bash
# Instead of multiple individual shell commands:
mv file1.txt file1.bak
mv file2.txt file2.bak
mv file3.txt file3.bak
# ...

# Use a Python one-liner or script:
python3 -c "
import pathlib
for f in pathlib.Path('.').glob('*.txt'):
    f.rename(f.with_suffix('.bak'))
"
```

**When to prefer Python scripts:**

- **Batch file operations** — rename, move, search-and-replace across files
- **Parsing or transforming structured data** — JSON, YAML, CSV
- **Code generation or templating tasks** — scaffolding boilerplate code
- **Any task requiring loops, conditionals, or string manipulation across multiple targets**

## Additional Use Cases

### Code Style Compliance Scanning

Scan the entire codebase to detect violations such as `!x` instead of `x == false`, plural variable names, or missing block ending comments. A single Python script using regex or AST parsing can check all files at once and produce a violation report.

```bash
python3 -c "
import pathlib, re
for f in pathlib.Path('.').rglob('*.go'):
    for i, line in enumerate(f.read_text().splitlines(), 1):
        if re.search(r'if\s+!', line):
            print(f'{f}:{i}: {line.strip()}')
"
```

### Proto / Sheet Change Impact Analysis

When `.proto` definitions or sheet schemas change, parse the files to automatically list which files across server, client, and gmtool are affected. Can also generate skeleton code for the required modifications.

### Test Data Generation

Batch-generate mock data conforming to project conventions (e.g., `helps.Date(2023, 2, 10)` patterns, typed structs) based on schemas, instead of writing each entry manually.

### Multi-File Synchronized Refactoring

When a message, enum, or interface is renamed, perform coordinated search-and-replace across model, view, router, i18n, and message files in a single execution, then verify consistency.

### Log and Output Parsing

After running tests or linters, parse the raw output to extract only failures and warnings into a concise summary, avoiding the need to read through large volumes of passing results.

### Excel / JSON Data Validation

Validate Sheeter-generated JSON files for format correctness, missing fields, or type mismatches without running the full build pipeline. Useful for quick pre-commit sanity checks.

# Commit / PR Conventions

## Commit Message Format

```text
<Type> | <Description in Traditional Chinese>
```

| Type | Usage |
| ---- | ----- |
| `Feature` | New feature |
| `Fix` | Bug fix |
| `Sheet` | Sheet/data table changes |
| `Message` | Message/proto updates |
| `UI` | UI adjustments |

**Examples:**

```text
Feature | 英雄準備攻擊的表演
Fix | 修正接關完成後副本的bug
Sheet | 調整英雄聖物坐騎相關能力數值
```

## Branch Naming

Format: `<account>/<feature-name>` (feature name in kebab-case English)

```text
mike/fight
xander_yuan/ride
yilin/fund
```

## Branch Strategy

| Branch | Purpose |
| ------ | ------- |
| `dev` | Development (main working branch) |
| `qa-xxx` | Internal testing |
| `preview` | Pre-release / review / staging |
| `main` | Production |

## PR Rules

- Feature branches → PR → `dev`
- PR description should summarize changes in Traditional Chinese
- When AI drafts a commit message or PR description, always follow the format above

# Context Length Self-Management

When operating in long conversations or processing large files, AI should proactively manage its context window to avoid missing instructions or losing important details.

## Guidelines

**Summarize proactively** — When a conversation exceeds roughly 20 exchanges, summarize the key decisions, current task state, and pending items before continuing. This prevents earlier instructions from being silently dropped.

**Process large files in segments** — When reading or modifying files over 500 lines, work in focused segments rather than loading everything at once. State which segment is being processed and track progress explicitly.

**Restate constraints before execution** — Before performing a complex multi-step task, briefly restate the relevant coding style rules and project constraints that apply. This serves as a self-check against drifting from the guidelines.

**Flag context pressure** — If the conversation has accumulated substantial content (multiple large code blocks, lengthy discussions), proactively inform the user and suggest starting a new conversation or summarizing the current state to a file for continuation.
