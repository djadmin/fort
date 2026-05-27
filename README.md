# fort

**Audit, fix, and prove macOS endpoint security from one command.**

[![CI](https://github.com/djadmin/fort/actions/workflows/ci.yml/badge.svg)](https://github.com/djadmin/fort/actions)
[![Release](https://img.shields.io/github/v/release/djadmin/fort)](https://github.com/djadmin/fort/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![macOS 12+](https://img.shields.io/badge/macOS-12%2B-lightgrey)](https://github.com/djadmin/fort)

`fort` checks your Mac's security settings, fixes the ones it can, and produces an auditor-ready HTML report — no MDM, no agent, no signup. Built for startups preparing for SOC 2, consultants doing client readiness reviews, and BYOD teams that need evidence without full device enrollment.

**[djadmin.github.io/fort](https://djadmin.github.io/fort)**

```
$ fort

  fort v0.1.0  —  alice-mbp (macOS 15.5)
  ───────────────────────────────────────────────────────────────────

  ✓  Password manager         1Password
  ✓  Disk encryption          on
  ✗  Screen lock              off                         expected: immediate
  ~  Antivirus / EDR          XProtect only (built-in)   expected: third-party AV/EDR
  ✓  Application firewall     on
  ✓  Gatekeeper               enabled
  ✓  System integrity (SIP)   enabled
  ✓  Remote login (SSH)       off
  ✗  Local admin rights       admin                      expected: standard user
  ✓  Guest account            disabled
  ✓  Automatic login          disabled
  ✓  Sharing services         all off
  ✓  AirDrop                  Off
  ✗  Automatic OS updates     off                        expected: on
  ✓  OS patch status          15.5

  ───────────────────────────────────────────────────────────────────
  Score: 11/15  (11 pass, 3 fail, 1 warn)

  Run fort --fix to remediate fixable issues.
```

## Install

```bash
# macOS (Apple Silicon + Intel)
curl -fsSL https://github.com/djadmin/fort/releases/latest/download/fort_darwin_all.tar.gz | tar xz && sudo mv fort /usr/local/bin/

# Go
go install github.com/djadmin/fort/cmd/fort@latest
```

## Usage

```bash
fort                # audit your Mac
fort --fix          # audit + auto-remediate fixable issues
fort --dry-run      # show exactly what --fix would change, without applying it
fort --json         # structured JSON output for scripts and dashboards
fort --report       # write fort-report-YYYY-MM-DD.html (auditor-ready, print to PDF)
```

Exit codes: `0` all pass · `1` any fail · `2` any warn

## What It Checks

15 macOS checks across five groups, each mapped to SOC 2, ISO 27001, NIST CSF, and CIS v8:

| Group | Checks |
|-------|--------|
| Core security | password manager, FileVault, screen lock, antivirus / EDR |
| System hardening | firewall, Gatekeeper, SIP, SSH |
| Access controls | local admin rights, guest account, automatic login |
| Exposure reduction | sharing services, AirDrop |
| Patching | automatic OS updates, OS patch status |

## Evidence report

`fort --report` writes a self-contained HTML file — machine identity, serial, OS version, timestamp, per-check results, and framework control references (SOC 2 CC6.x / CC7.x, ISO 27001 A.8.x, NIST CSF, CIS v8). Open in any browser and print to PDF. No server, no upload.

## JSON output

`fort --json` emits a stable payload for automation:

```json
{
  "tool": "fort",
  "version": "0.1.0",
  "hostname": "alice-mbp",
  "serial": "C02X...",
  "os_version": "15.5",
  "timestamp": "2026-05-21T10:30:00Z",
  "summary": { "total": 15, "pass": 11, "fail": 3, "warn": 1, "score": "11/15" },
  "policies": [
    {
      "id": "filevault",
      "name": "Disk encryption",
      "status": "pass",
      "current": "on",
      "expected": "on",
      "fixable": false,
      "fixed": false,
      "frameworks": { "SOC 2": ["CC6.1", "CC6.7"], "ISO 27001": ["A.8.3", "A.8.24"], ... }
    }
  ]
}
```

## When to use fort

fort is not an MDM replacement. It is useful when you need to:

- understand a Mac's current security posture quickly
- remediate obvious gaps before an audit or client review
- produce timestamped evidence without enrolling devices in a SaaS platform

For full MDM (remote wipe, app deployment, profile management), use Jamf, Kandji, or Intune.

## Contributing

PRs are welcome. To add a check:

1. Create `internal/checks/yourcheck_darwin.go` — implement the `Check` interface
2. Register it in `internal/checks/registry_darwin.go`
3. Add framework mappings in `internal/checks/frameworks.go`
4. Run `go test ./...` — existing tests enforce interface contracts

New checks should use documented macOS APIs, work on macOS 12+, and have clear pass/fail criteria.

## License

[MIT](LICENSE)
