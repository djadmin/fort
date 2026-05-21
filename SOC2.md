# fort × SOC 2

What SOC 2 requires on endpoints, what fort covers, what's missing, and how we compare to similar tools.

---

## SOC 2 in plain English

SOC 2 is a third-party security audit. A CPA firm reviews whether your company's security controls actually work — not just that you wrote a policy, but that the policy is enforced. Customers (especially enterprise) demand it before sharing data with you.

**Type I vs. Type II — the distinction that matters:**

| | Type I | Type II |
|--|--|--|
| Proves | Controls *exist* today | Controls *worked* for 6–12 months |
| Evidence | One example per control | Sampled across the whole period |
| Timeline | 3–6 months total | 9–15 months total |
| What buyers want | Early-stage, first audit | Almost always required eventually |
| Fort's role | HTML report covers this well | Needs fleet aggregation (v0.4) |

---

## The criteria that drive endpoint requirements

SOC 2 is organized into Trust Services Criteria (TSC). Only Security is mandatory. Within Security, the "Common Criteria" (CC) series. The ones that apply to laptops and endpoints:

```
CC6.1   Who can access the system
        → Screen lock, MFA, no shared accounts, logical access controls

CC6.3   Least privilege
        → Users should not run as local admin; access matches job function

CC6.6   Boundary protection
        → Firewall on, SSH off, AirDrop restricted, sharing services disabled, VPN for remote

CC6.7   Data in transit and at rest
        → Full-disk encryption (FileVault), USB/removable media controls

CC6.8   Malicious software prevention
        → AV/EDR deployed, Gatekeeper on, OS patched, SIP enabled, software allowlisting

CC7.1   Detection and monitoring
        → Logs forwarded, EDR telemetry flowing, patch compliance tracked

CC7.3   Recovery (if Availability is in scope)
        → Encrypted backups (Time Machine), recovery procedures tested
```

---

## What auditors actually ask for (evidence guide)

| Control | What auditors accept | What they reject |
|--|--|--|
| FileVault | MDM console export showing per-device status; `fdesetup status` output | "We told employees to turn it on" |
| Screen lock | MDM config profile showing enforcement; fort HTML report per machine | No enforcement evidence |
| AV/EDR | EDR console showing all devices enrolled + agent version | XProtect-only (usually) for Type II |
| Firewall | MDM screenshot or CLI output | Policy doc without enforcement |
| OS patches | MDM fleet report showing OS versions; patch SLA policy | Updates "on" without version check |
| Password manager | Software license list + MDM inventory showing installed | Honor system |
| SSH off | CLI output; MDM config | Nothing |
| Admin rights | MDM report showing standard vs. admin account split | Self-reported |

**Fort's HTML report (`fort --report`) is good evidence for Type I.** For Type II, auditors need to see this data for all machines, across the full observation period — that's the fleet dashboard.

---

## Fort's current SOC 2 coverage

### ✅ Covered (8 checks)

| Check | Criterion | Auto-fix |
|--|--|--|
| Password manager installed | CC6.1 | — |
| FileVault (disk encryption) | CC6.7, CC6.1 | — (needs reboot) |
| Screen lock | CC6.1, CC6.6 | ✅ |
| Antivirus / EDR presence | CC6.8, CC7.1 | — |
| Application firewall | CC6.6 | ✅ (sudo) |
| Gatekeeper | CC6.8 | ✅ (sudo) |
| Remote login (SSH) off | CC6.6 | ✅ (sudo) |
| Automatic OS updates | CC6.8 | ✅ (sudo) |

### ❌ Missing — high priority (next sprint, target: v0.1.1)

| Check | Criterion | Why it matters |
|--|--|--|
| Local admin rights check | CC6.3, CC6.8 | Most commonly cited audit finding. Users running as admin can install malware, bypass controls. Vanta, Drata, Pareto all check this. |
| Guest account disabled | CC6.1 | Physical access to a Mac with guest enabled bypasses logical access controls. |
| System Integrity Protection (SIP) | CC6.8 | SIP prevents root-level modification of system files. Disabling it is a red flag. Gatekeeper and XProtect depend on SIP. |
| Sharing services off | CC6.6 | File Sharing, Printer Sharing, Internet Sharing, Remote Management — each is an attack surface. Pareto checks all 5. |
| Automatic login disabled | CC6.1 | Auto-login machine = physical access bypasses all auth controls. |
| AirDrop restricted | CC6.6 | Unintended data exfiltration vector. Should be set to "Contacts Only" or "Off". |
| Current OS version (not just updates on) | CC6.8 | Having updates "on" doesn't mean patches are applied. Auditors check the specific OS version against current release. |
| Secure Boot | CC6.1 | Prevents boot-level persistence and unauthorized OS booting. |

### ⚠️ Partial coverage

| Check | Issue |
|--|--|
| Screen lock timeout | We check on/off + immediate delay. Auditors also want to see the idle timeout (should be ≤ 5–10 min for laptops). |
| EDR "active" vs. "installed" | We detect presence (app/process). We don't verify the agent is connected to a management console or that telemetry is flowing. |
| OS updates | We check the setting is on. We should also check the current OS version against latest release. |

### 🏗️ Structural gap — Type II

Fort produces per-machine, point-in-time evidence. For SOC 2 Type II, auditors need:
- Evidence from **all** managed endpoints (not just one)
- Evidence maintained **throughout** the 6–12 month observation period
- **Drift detection** (was a machine ever non-compliant during the period?)

This requires the fleet dashboard (v0.4 in BACKLOG.md). The JSON output (`fort --json`) is the foundation — machines can POST results to a collector.

---

## Competitive analysis

### Direct competitors

**Pareto Security** (paretoapp.com) — closest comparable
- macOS menubar app; 32 checks; runs automatically in background
- Pricing: $17 personal (one-time) / subscription for teams / $150/mo enterprise
- Vanta integration (can push evidence directly to Vanta)
- Open-source desktop version available
- **No auto-remediation** — reports only, never fixes
- No JSON output, no CLI, no scripting
- All 32 checks: FileVault, firewall, Gatekeeper, screen lock, password manager, SSH off, AirDrop, all sharing services, secure boot, SIP, no admin account, no guest account, Time Machine encrypted backup, uptime < 14 days (patch applied), SSH key strength, terminal secure keyboard, WiFi security, and more

**Where fort beats Pareto:**
- Auto-remediation (fort fixes; Pareto only reports)
- JSON output (fort pipes to fleet, SIEM, Slack; Pareto doesn't)
- CLI / scriptable (fort works in MDM push, CI/CD; Pareto is a GUI app)
- Evidence report (fort writes auditor-ready HTML; Pareto has no export)
- Free

**Where Pareto beats fort (for now):**
- 32 checks vs. 8
- Vanta/Drata integration
- Team dashboard
- Runs continuously in background (menubar icon)

---

**Lynis** — open source, CLI, macOS + Linux
- 200+ checks across 38 categories; Hardening Index score
- Server-focused: kernel params, NFS, SNMP, webservers, SSH server hardening
- Not SOC 2 mapped; not laptop-focused
- No auto-remediation; no JSON; no evidence report
- Good for security hardening; not positioned as compliance evidence

**CIS-CAT** (Center for Internet Security)
- Java-based scanner against CIS Benchmarks (hundreds of controls)
- Maps to PCI DSS, NIST 800-53, HIPAA, ISO 27001, DISA STIG
- Generates HTML/CSV compliance reports
- Requires CIS SecureSuite membership (paid enterprise)
- No auto-remediation; not lightweight
- Overkill for startups; designed for enterprise IT teams

**Kolide / 1Password Device Trust**
- osquery-based agent; 100+ pre-written checks; custom check support
- Conditional access: blocks SSO login if device fails checks
- Integrates with Okta, Vanta, Drata
- Enterprise pricing; not open source
- Fort's `--json` output is designed to eventually integrate with systems like this

**NIST macOS Security Compliance Project (mSCP)**
- Official NIST project; maps to NIST SP 800-53, DISA STIG
- CLI generates MDM configuration profiles, compliance scripts, guidance docs
- Government/federal grade — much heavier than SOC 2 typically requires
- Open source; available at github.com/usnistgov/macos_security

---

### Summary comparison

| | fort | Pareto | Lynis | CIS-CAT | Kolide |
|--|--|--|--|--|--|
| macOS endpoint checks | 8 (→ 15) | 32 | Partial | 100+ | 100+ |
| Auto-remediation | ✅ | ❌ | ❌ | ❌ | ❌ |
| CLI / scriptable | ✅ | ❌ | ✅ | ✅ (Java) | ❌ |
| JSON output | ✅ | ❌ | ❌ | CSV | API |
| Evidence report | ✅ HTML | ❌ | ❌ | ✅ HTML | ✅ dashboard |
| SOC 2 mapped | Partial | ❌ | ❌ | ✅ | ✅ |
| Fleet view | 🔜 v0.4 | ✅ paid | ❌ | ❌ | ✅ enterprise |
| Open source | ✅ | ✅ | ✅ | ❌ | ❌ |
| Price | Free | $17–$150/mo | Free | Paid | Enterprise |
| Windows / Linux | 🔜 v0.3 | ❌ | ✅ | ✅ | ✅ |

---

## What to build next (prioritized by SOC 2 impact)

### Sprint 1 — get to 15 checks (closes most SOC 2 endpoint gaps)

1. `localadmin` — detect if user is running as local admin (CC6.3)
2. `guestaccount` — guest account disabled (CC6.1)
3. `sip` — System Integrity Protection enabled (CC6.8)
4. `filesharing` — File Sharing + Printer Sharing + Remote Management + Internet Sharing all off (CC6.6)
5. `autologin` — automatic login disabled (CC6.1)
6. `airdrop` — AirDrop set to Contacts Only or Off (CC6.6)
7. `osversion` — check macOS version against current release, not just updates setting (CC6.8)

### Sprint 2 — evidence quality

8. Screen lock timeout value (currently checks on/off; add idle timeout threshold check)
9. SOC 2 control mapping in `fort --report` (CC numbers next to each check)
10. `fort --report --format pdf` (via headless Chrome)

### Sprint 3 — fleet + Type II

11. Fleet dashboard (v0.4) — POST JSON results to a collector, show all machines, flag drift
12. Vanta/Drata webhook integration — push results automatically on scan
