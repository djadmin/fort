# fort

**Audit, fix, and prove macOS endpoint security from one command.**

`fort` is a macOS CLI for teams that need fast, auditable endpoint checks without standing up MDM first. It inspects local security settings, can remediate the fixable ones, and can emit both machine-readable JSON and an auditor-friendly HTML report.

Built for startups preparing for SOC 2, consultants running client readiness reviews, and BYOD environments that need evidence without full device enrollment.

```text
$ fort

  fort v0.1.0  —  alice-mbp (macOS 15.5)
  ───────────────────────────────────────────────────────────────────

  ✓  Password manager         1Password
  ✓  Disk encryption          on
  ✗  Screen lock              off                         expected: immediate
  ~  Antivirus / EDR          XProtect only              expected: third-party AV/EDR
  ✓  Application firewall     on
  ✓  Gatekeeper               enabled
  ✓  Remote login (SSH)       off
  ...

  ───────────────────────────────────────────────────────────────────
  Score: 11/15  (11 pass, 2 fail, 2 warn)

  Run fort --fix to remediate fixable issues.
```

## Status

- macOS 12+ today
- Single local binary, no agent, no signup
- 15 built-in checks mapped to SOC 2, ISO 27001, NIST CSF, and CIS v8

## Install

```bash
go install github.com/djadmin/fort/cmd/fort@latest
```

Pre-built binaries and a Homebrew tap are coming with the first tagged release.

## Usage

```bash
fort                # audit your Mac
fort --fix          # audit + auto-remediate fixable issues
fort --dry-run      # print what --fix would change
fort --json         # structured output for automation
fort --report       # write fort-report-YYYY-MM-DD.html
```

## What It Checks

`fort` currently ships 15 macOS checks across five groups:

| Area | Checks |
|------|--------|
| Core security | password manager, FileVault, screen lock, antivirus / EDR |
| System hardening | firewall, Gatekeeper, SIP, SSH |
| Access controls | local admin rights, guest account, automatic login |
| Exposure reduction | sharing services, AirDrop |
| Patching | automatic OS updates, pending OS updates |

Each result includes framework mappings in both JSON and HTML report output.

## Outputs

`fort --json` produces a stable JSON payload designed for scripts, collectors, and dashboards:

```json
{
  "tool": "fort",
  "version": "0.1.0",
  "hostname": "alice-mbp",
  "serial": "C02X...",
  "os_version": "15.5",
  "timestamp": "2026-05-21T10:30:00Z",
  "summary": { "total": 15, "pass": 11, "fail": 2, "warn": 2, "score": "11/15" },
  "policies": [...]
}
```

`fort --report` writes a self-contained HTML evidence report with machine identity, timestamp, per-check status, fix markers, and framework references. The file has no server dependency and can be opened locally or printed to PDF.

## Positioning

`fort` is not an MDM replacement. It is the fast path to:

- understand a Mac's current security posture
- remediate obvious gaps
- produce evidence an auditor or client can review

Compared with larger tools, the value is speed, transparency, and a workflow that starts with one command instead of enrollment and SaaS setup.

## Repository Layout

| Path | Purpose |
|------|---------|
| `cmd/fort` | CLI entrypoint, output, and report generation |
| `internal/checks` | macOS checks and framework mappings |
| `cmd/landing` | optional landing-page waitlist server |
| `landing` | static marketing site assets |

## Contributing

PRs are welcome.

1. Add a new check under `internal/checks`.
2. Register it in `internal/checks/registry_darwin.go`.
3. Add framework mappings in `internal/checks/frameworks.go`.
4. Run `go test ./...`.

New checks should use documented macOS interfaces, work on supported macOS versions, and have clear pass/fail semantics.

## License

[MIT](LICENSE)
