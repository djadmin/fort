# fort

**One command. Every Mac. SOC 2 ready.**

fort audits your Mac's security settings, fixes what it finds, and produces a timestamped compliance report — no MDM, no agent, no signup.

Built for: startups doing their first SOC 2, security consultants, BYOD teams that need proof without the overhead of Jamf or Kandji.

```
$ fort

  fort v0.1.0  —  alice-mbp (macOS 15.5)
  ───────────────────────────────────────────────────────

  ✓  Password manager       1Password
  ✓  Disk encryption        on
  ✗  Screen lock            off               expected: immediate
  ~  Antivirus / EDR        XProtect only     expected: third-party AV/EDR
  ✓  Application firewall   on
  ✓  Gatekeeper             enabled
  ✓  Remote login (SSH)     off
  ✗  Automatic OS updates   off               expected: on

  ───────────────────────────────────────────────────────
  Score: 5/8  (5 pass, 2 fail, 1 warn)

  Run fort --fix to remediate fixable issues.
```

## Install

```bash
# Homebrew (coming soon)
brew install djadmin/tap/fort

# Manual
curl -fsSL https://github.com/djadmin/fort/releases/latest/download/fort-darwin-arm64 -o fort
chmod +x fort && sudo mv fort /usr/local/bin/fort
```

## Usage

```bash
fort                # audit your Mac
fort --fix          # audit + auto-remediate fixable issues (some require sudo)
fort --dry-run      # show exactly what --fix would change, without touching anything
fort --json         # structured JSON output (pipe to dashboard, SIEM, Slack)
fort --report       # audit + write fort-report-YYYY-MM-DD.html (auditor-ready evidence)
```

## Checks

| # | Check | How it detects | Auto-fix |
|---|-------|---------------|----------|
| 1 | Password manager | Scans `/Applications` for known password managers | — |
| 2 | Disk encryption (FileVault) | `fdesetup status` | — (needs reboot) |
| 3 | Screen lock | `defaults read com.apple.screensaver` | ✅ |
| 4 | Antivirus / EDR | App presence + running processes | — |
| 5 | Application firewall | `socketfilterfw --getglobalstate` | ✅ (sudo) |
| 6 | Gatekeeper | `spctl --status` | ✅ (sudo) |
| 7 | Remote login (SSH) | `launchctl print system/com.openssh.sshd` | ✅ (sudo) |
| 8 | Automatic OS updates | `defaults read .../SoftwareUpdate` | ✅ (sudo) |

Every check uses stable, documented macOS APIs. No private frameworks. macOS 12+.

## Evidence report

`fort --report` writes a self-contained HTML file you can open in a browser and print to PDF — no external dependencies, no server.

Includes: machine identity, OS version, serial number, timestamp, per-check pass/fail/warn, what was auto-fixed.

Designed to satisfy auditor requests for "evidence of endpoint security controls."

## JSON output

```bash
fort --json
```

```json
{
  "tool": "fort",
  "version": "0.1.0",
  "hostname": "alice-mbp",
  "serial": "C02X...",
  "os_version": "15.5",
  "timestamp": "2026-05-21T10:30:00Z",
  "summary": { "total": 8, "pass": 6, "fail": 1, "warn": 1, "score": "6/8" },
  "policies": [...]
}
```

Pipe anywhere — dashboard, SIEM, Slack webhook, your own collector.

## Why not just use...

| | fort | Jamf/Kandji | Vanta/Drata | Lynis | osquery |
|--|------|-------------|-------------|-------|---------|
| Zero setup | ✅ | ❌ MDM enrollment | ❌ agent + SaaS | ✅ | ❌ |
| Auto-remediation | ✅ | ✅ | ❌ | ❌ | ❌ |
| Evidence report | ✅ | ✅ | ✅ | ❌ | ❌ |
| Open source | ✅ | ❌ | ❌ | ✅ | ✅ |
| Free | ✅ | ❌ $$$  | ❌ $$$ | ✅ | ✅ |
| Windows/Linux | 🔜 v0.3 | ✅ | ✅ | ✅ | ✅ |

fort is not trying to replace Jamf. It's what you use before you need Jamf — or when a contractor needs to prove compliance without enrolling in your MDM.

## Fleet use

```bash
# Run on each machine (via MDM push, Munki, or SSH loop):
fort --json > /tmp/fort-$(hostname).json

# POST to your collector:
curl -s -X POST https://your-dashboard.com/api/report \
  -H "Content-Type: application/json" \
  -d @/tmp/fort-$(hostname).json
```

## Philosophy

1. **Transparent.** Every check is a CLI command you can run yourself. The source is the documentation.
2. **Opinionated.** Eight checks that matter, reliably detected. Not 200 checks that are 60% guesswork.
3. **Honest.** Checks that can't auto-fix say so. Checks that need sudo say so. No false confidence.

## Contributing

PRs welcome. To add a check:
1. Create `internal/checks/yourcheck_darwin.go` implementing the `Check` interface
2. Add it to `All()` in `internal/checks/registry_darwin.go`
3. Write a test — `go test ./...` must pass

New checks must: use documented macOS APIs, work on macOS 12+, have clear pass/fail criteria.

## License

MIT
