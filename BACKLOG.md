# fort — backlog

## Status: v0.1.0 — shipped locally, not yet on GitHub

---

## Done

- [x] 8 macOS checks: password manager, FileVault, screen lock, AV/EDR, firewall, Gatekeeper, SSH, OS auto-updates
- [x] `--fix` with per-check auto-remediation
- [x] `--dry-run` showing exact commands before applying
- [x] `--json` structured output (stable fleet contract)
- [x] `--report` HTML evidence artifact (print-to-PDF ready)
- [x] Score display with pass/fail/warn counts
- [x] Go build, `make install` to `~/.local/bin`

---

## Now — launch prep (this week)

- [ ] **Git init + push to github.com/djadmin/fort** (public, MIT)
- [ ] **Code sign + notarize binary** — required for `--fix` trust. Without this, macOS may block execution for downloaded binaries. Use Apple Developer account.
- [ ] **GoReleaser config** — multi-arch builds (darwin/amd64, darwin/arm64), GitHub releases, checksums, SBOM
- [ ] **Homebrew tap** — `brew install djadmin/tap/fort`
- [ ] **Landing page** — single HTML page (no framework): headline "Get every Mac ready for SOC 2 in 10 minutes", terminal demo gif, example HTML report screenshot, email waitlist for team features

---

## Next — v0.2 (evidence + policy)

- [ ] **`fort.yaml` policy file** — org-defined baseline, checked against at runtime
  ```yaml
  password_manager:
    allowed: [1Password, Bitwarden]
  edr:
    allowed: [CrowdStrike Falcon, SentinelOne]
  screen_lock: immediate
  disk_encryption: required
  ```
- [ ] **PDF report** — via headless Chrome (`chromium --headless --print-to-pdf`) or wkhtmltopdf. Needed for auditors who want PDF, not HTML.
- [ ] **SOC 2 control mapping** — each check maps to CC6.x / CC7.x controls. Shown in report.
- [ ] **ISO 27001 control mapping** — A.8.x mappings. Shown in report alongside SOC 2.
- [ ] **`fort --check <id>`** — run a single check by ID
- [ ] **Scheduled scan + report** — cron-friendly: `fort --json --report --quiet` exits 0/1 on pass/fail
- [ ] **Exit code semantics** — 0 = all pass, 1 = any fail, 2 = any warn (for CI/CD gates)

---

## Then — v0.3 (multi-OS)

- [ ] **Windows agent** — PowerShell checks: BitLocker, Defender, Windows Firewall, UAC, RDP, guest account, lock screen, Windows Update. Same JSON contract.
- [ ] **Linux agent** — UFW/iptables, SSH root login, SSH password auth, unattended upgrades, fail2ban, SELinux/AppArmor. Same JSON contract.
- [ ] Both compile from same Go codebase via `_windows.go` / `_linux.go` build tags

---

## Later — v0.4 fleet dashboard (paid)

- [ ] **Collector API** — POST endpoint that accepts `fort --json` output, stores per-machine results
- [ ] **Fleet view** — web UI: all machines, compliance %, per-machine drill-down, last-seen
- [ ] **Drift detection** — alert when a passing machine starts failing
- [ ] **Slack / email alerts** — notify on new device or compliance regression
- [ ] **Team enrollment** — `fort enroll --org acme --token xxx` sets org + auto-posts results

---

## Ideas parking lot

- **Shareable score badge** — `fort --badge` outputs a URL with score. Viral distribution in Slack/all-hands.
- **GitHub Actions integration** — `fort --ci` exits non-zero if score below threshold. Block deploys on non-compliant machines.
- **Browser extension checks** — detect 1Password / Bitwarden browser extensions (harder, but useful)
- **Login item audit** — flag suspicious startup items
- **VPN check** — detect if a corporate VPN client is installed
- **`fort verify <report.html>`** — verify a report's integrity (HMAC or checksum)
- **White-label reports** — for consultants: custom logo, client name in report header
- **Consultant workspace** — one org per client, shared policy templates, batch report export

---

## Positioning

**Who pays:**
1. Startups (10-100 employees) preparing for SOC 2 / ISO 27001 — need proof fast, don't have IT team
2. vCISOs / compliance consultants — need before/after reports for clients, repeatable workflow
3. BYOD teams — contractors need to self-attest without enrolling in MDM

**Who doesn't pay (but uses it):**
- Individual developers — open-source users, GitHub stars, organic distribution

**Pricing hypothesis (validate before building dashboard):**
- CLI: free, forever
- Team (dashboard): $4-6/device/mo, min 5 devices
- Consultant workspace: $199-299/mo for unlimited client orgs + white-label reports
- One-time audit pack: $99 per company for 30-day compliance sprint (landing page CTA)

**Positioning that works:**
- "Get every Mac ready for SOC 2 in 10 minutes" (specific, urgent, clear buyer)
- "Fix and prove endpoint security before your audit" (action-oriented)
- NOT: "CLI security checker" (sounds like a free dev tool)

**Differentiation:**
- Remediation (Lynis/osquery don't fix, just audit)
- Evidence report (what auditors actually ask for)
- Transparent (open source, one binary, read the source)
- Fast to start (no MDM enrollment, no SaaS signup, `curl | install`)

---

## Competitive watch

- **1Password Device Trust (Kolide)** — acquired for ~$100M, enterprise-heavy, Slack-based. Gap: too heavy for < 50 person teams.
- **Vanta / Drata device monitoring** — detect but don't remediate. Tied to their compliance platform.
- **Fleet** — open-source device management, powerful but complex. Gap: no remediation, high setup cost.
- **Pareto Security** — macOS-only GUI app, $5/mo personal. No fleet view, no evidence export.
- **CIS-CAT** — benchmarking tool, no remediation, complex licensing.
