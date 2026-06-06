# Changelog

All notable changes to fort are documented here.

---

## v0.2.0 — 2026-05-28

**Homebrew tap + distribution polish**

fort is now installable in one line via Homebrew — the recommended path for most users:

```bash
brew install djadmin/tap/fort
```

- Added Homebrew tap (`djadmin/tap/fort`) as the primary install method
- Three supported install paths: Homebrew, pre-built binary, `go install`
- Landing page at [djadmin.github.io/fort](https://djadmin.github.io/fort) with screenshots and waitlist
- README rewritten: side-by-side audit/fix screenshots, tighter copy, clearer individual vs. compliance positioning

---

## v0.1.1 — 2026-05-27

**Interactive fix workflow**

The `--fix` flag now walks you through each remediation interactively before changing anything.

- Interactive confirmation prompt before applying any fix — lists exactly what will change
- Selective apply: choose which fixes to run rather than accepting all
- `--yes` / `-y` flag for non-interactive use (automation, CI, MDM scripts)
- Improved test coverage across checks, output, and report

---

## v0.1.0 — 2026-05-27

**Initial release**

fort's first public release. Runs 15 security checks on macOS, maps every finding to SOC 2 / ISO 27001 / NIST CSF / CIS v8, and produces a self-contained HTML evidence report.

- **15 checks** across disk encryption, screen lock, firewall, Gatekeeper, SIP, SSH, AirDrop, guest account, auto-login, automatic updates, OS version, antivirus, password manager, local admin, and sharing services
- `fort --fix` to remediate fixable gaps (with confirmation)
- `fort --dry-run` to preview what a fix would change
- `fort --json` for machine-readable output (CI / scripting)
- `fort --report` to write an auditor-ready, self-contained HTML report
- Exit codes: `0` = all pass, `1` = any fail, `2` = any warn
- Framework mappings on every check: SOC 2 CC controls, ISO 27001 A-controls, NIST CSF subcategories, CIS v8 safeguards
- No agent, no signup, no MDM — single static binary, zero runtime dependencies
