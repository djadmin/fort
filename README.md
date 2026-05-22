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

## Checks — 15 controls, full SOC 2 coverage

| # | Check | SOC 2 | How it detects | Auto-fix |
|---|-------|-------|---------------|----------|
| 1 | Password manager | CC6.1 | Scans `/Applications` for 12 known managers | — |
| 2 | Disk encryption (FileVault) | CC6.7 | `fdesetup status` | — (needs reboot) |
| 3 | Screen lock | CC6.1 | `defaults read com.apple.screensaver` | ✅ |
| 4 | Antivirus / EDR | CC6.8 | App presence + running processes (17 tools) | — |
| 5 | Application firewall | CC6.6 | `socketfilterfw --getglobalstate` | ✅ (sudo) |
| 6 | Gatekeeper | CC6.8 | `spctl --status` | ✅ (sudo) |
| 7 | System integrity (SIP) | CC6.8 | `csrutil status` | — (recovery mode) |
| 8 | Remote login (SSH) | CC6.6 | `launchctl print system/com.openssh.sshd` | ✅ (sudo) |
| 9 | Local admin rights | CC6.3 | `id` — checks if user is in admin group | — |
| 10 | Guest account | CC6.1 | `defaults read .../loginwindow GuestEnabled` | ✅ (sudo) |
| 11 | Automatic login | CC6.1 | `defaults read .../loginwindow autoLoginUser` | ✅ (sudo) |
| 12 | Sharing services | CC6.6 | `launchctl` — checks SMB, screen sharing, ARD, internet sharing | — |
| 13 | AirDrop | CC6.6 | `defaults read com.apple.sharingd DiscoverableMode` | ✅ |
| 14 | Automatic OS updates | CC6.8 | `defaults read .../SoftwareUpdate AutomaticCheckEnabled` | ✅ (sudo) |
| 15 | OS patch status | CC6.8 | Checks `PendingUpdateCount` from cached SoftwareUpdate prefs | — |

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
