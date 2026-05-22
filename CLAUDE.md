# fort — project context

Endpoint security CLI for macOS (Windows/Linux in v0.3). Open-core: CLI is free/OSS, fleet dashboard is paid.

**Target**: Startups doing SOC 2, vCISOs, BYOD teams. Not a Jamf replacement.

## Build

```bash
make build       # → ./fort
make test        # go test ./...
make install     # → ~/.local/bin/fort (GOBIN set in ~/.zshrc)
```

Go is managed via mise. Commands must be prefixed with `mise exec --` or use `make` (Makefile handles this).

## Architecture

```
cmd/fort/
  main.go       Cobra CLI, flag parsing, run() orchestration, sysinfo helpers
  output.go     Human-readable terminal output + JSON
  report.go     HTML evidence report (self-contained, print-to-PDF)

internal/checks/
  check.go              Check interface + Result/Status types
  registry_darwin.go    All() — ordered check list for macOS
  registry_other.go     stub returning nil for non-darwin
  *_darwin.go           Per-check implementations (build tag: darwin only)
  checks_test.go        Interface compliance + smoke tests for all checks
```

## Adding a check (macOS)

1. Create `internal/checks/yourcheck_darwin.go`
2. Implement `Check` interface: `ID()`, `Name()`, `Run()`, `Fixable()`, `Fix()`, `FixDescription()`
3. Non-fixable: return `false`/`nil`/`""` for Fixable/Fix/FixDescription
4. **Every `Run()` branch must set `Fixable: c.Fixable()`** — `TestFixableConsistency` enforces this
5. Add to `All()` in `registry_darwin.go` and update `TestRegistryCount` expected count
6. Add framework mappings in `internal/checks/frameworks.go` — `TestAllChecksHaveFrameworkMappings` enforces this
7. Run `make test` — all 20+ tests must pass

## Testing rules

**Every feature and check must have test coverage. No exceptions.**

```
internal/checks/checks_test.go   — check interface contracts, registry, framework mappings
cmd/fort/output_test.go          — tally(), anyFixable(), toJSONPolicies()
cmd/fort/report_test.go          — writeReport() HTML output, all key elements
```

**What tests enforce (will break CI if violated):**
- `TestRegistryCount` — update the count when adding/removing a check
- `TestFixableConsistency` — every `Run()` branch must set `Fixable: c.Fixable()`
- `TestFixableChecksHaveDescription` — fixable checks must have `FixDescription()`
- `TestAllChecksHaveFrameworkMappings` — every check needs framework mappings
- `TestKnownFrameworks` — every check maps to all four frameworks (SOC 2, ISO 27001, NIST CSF, CIS v8)
- `TestNoDuplicateIDs` — check IDs must be globally unique
- `TestCheckIDsAreLowercase` — IDs must be lowercase, no whitespace

**What we don't mock (yet):** OS commands (`defaults`, `launchctl`, etc.) run for real in tests.
This means tests verify the current machine state, not simulated states. When the Windows/Linux
agents are added, use build tags to run platform-specific tests only on that OS in CI.

**CI:** GitHub Actions runs `go test ./...` on every push to `main` and on every tag.
The workflow is `.github/workflows/ci.yml`. Keep it green at all times.

## Key decisions

- **Language**: Go. Single static binary, cross-compiles, no runtime dep on target machine.
- **Module**: `github.com/djadmin/fort`
- **CLI framework**: `github.com/spf13/cobra`
- **GOBIN**: `~/.local/bin` (set in `~/.zshrc` — already in PATH)
- **Checks requiring sudo**: Fix() returns error if run without sudo; caller prints error and continues
- **Platform isolation**: `_darwin.go` filename suffix is enough; `//go:build darwin` added for explicitness

## Flags

```
fort                  run all checks, human output
fort --fix            run + auto-remediate fixable issues
fort --dry-run        show what --fix would do, nothing changed
fort --json           structured JSON (stable contract for fleet use)
fort --report         write fort-report-YYYY-MM-DD.html
```

## Check status semantics

- `pass` — meets the requirement
- `fail` — does not meet, remediation needed
- `warn` — partial (e.g. XProtect present but no third-party AV) — show yellow, don't count as fail in score

## Score

`pass / total` shown as fraction. `warn` counts neither pass nor fail. Score color: green (all pass), yellow (any warn, no fail), red (any fail).
