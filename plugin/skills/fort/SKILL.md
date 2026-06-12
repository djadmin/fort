---
name: fort
description: >-
  Use when the user wants to check, understand, or harden the security of their
  Mac, or asks whether their laptop is secure and what to fix. Runs a full macOS
  security audit, explains each finding in plain language (what it is, why it
  matters, what changing it would affect), and applies only the fixes the user
  approves. Also produces SOC 2 / ISO 27001 / NIST / CIS evidence on request.
  Triggers on phrases like "is my Mac secure", "audit my laptop", "harden my
  Mac", "what security gaps do I have", "fix my firewall / screen lock /
  FileVault", "check my endpoint security".
tools:
  allowed: ["Bash", "Read"]
---

# fort — audit and harden your Mac's security

`fort` is a single Go binary that runs 16 deterministic security checks on a Mac,
shows the exact command and output behind each result, and can fix the ones that
are safe to automate. This skill drives that CLI directly over the shell.

Your job is not to re-implement the checks. `fort` is the source of truth for *what*
the posture is. You add the layer a binary can't: explain **what each finding means
and why it matters**, prioritize them, respect what the user has told you about how
they use the machine, and apply fixes selectively and only with consent.

## The one rule that matters most

**Auditing is always safe. Fixing changes the user's system settings — treat every
fix like running a command on their machine: preview it, get an explicit yes, then
apply. Never apply a fix the user hasn't specifically approved, and never "fix
everything" implicitly.** Default to read-only. When in doubt, audit and explain;
don't change anything.

## Is fort installed?

Run `fort --version`. If the command isn't found, tell the user to install it and
stop there — this skill orchestrates the binary, it doesn't replace it:

```
brew install djadmin/tap/fort
```

## How to drive it

| Goal | Command | Changes the system? |
|---|---|---|
| Full audit, human-readable | `fort` | No |
| Full audit, structured data to reason over | `fort --json` | No |
| Audit specific checks only | `fort --only firewall,osupdates --json` | No |
| Preview the exact fix command for a check | `fort --only firewall --dry-run` | No |
| Apply a fix, with fort's own confirm prompt | `fort --only firewall --fix` | **Yes** |
| Apply a fix non-interactively (after the user approved) | `fort --only firewall --fix --yes` | **Yes** |
| Write an HTML evidence report | `fort --report` | No (writes `fort-report.html`) |

Notes:
- `fort` exits non-zero when any check fails. That's expected, not an error — read
  the output, don't treat a non-zero exit as a crash.
- Always prefer `--json` when you need to reason over results; it gives you per-check
  `status`, `current` vs `expected`, `fixable`, the `evidence` (the literal command
  fort ran and its output), and `frameworks` mappings.
- Scope fixes with `--only <id>`. A bare `fort --fix` would prompt for everything;
  you want to apply exactly the checks the user agreed to.

## Workflow

### 1. Audit first

Run `fort --json`. Lead with the score and a one-line verdict, then walk the findings
in priority order — failures that expose the machine first, lower-risk items after.
Don't dump JSON at the user.

### 2. Explain so they can actually see what's going on

This is the point of the skill. For each finding that matters, tell the user:
- **What it is** in plain language ("FileVault is the disk encryption that protects
  your data if the laptop is lost or stolen").
- **Where it stands** — `current` vs `expected`, and the `evidence` (the real command
  fort ran and what it returned). Showing the evidence is what makes this trustworthy
  instead of a black box.
- **Why it matters** — the concrete risk, not a compliance label.
- **What fixing it would affect** — anything that might break or change for the user.

See `reference/checks.md` for a per-check explainer (what it is, why it matters,
what to watch out for when changing it). Read it when you need depth on a finding or
the user asks "what happens if I change this?"

### 3. Respect the user's actual setup

A finding being "fail" doesn't mean blindly fix it. If the user says they use SSH to
reach this machine, don't push to disable remote login. If they need an admin account
for daily work, explain the trade-off rather than insisting. Prioritize against how
*they* use the Mac, not a generic checklist.

### 4. Preview before changing anything

When the user wants to fix something, run `fort --only <id> --dry-run` and show them
the exact command fort would run (`would run: sudo …`). Let them see it before it
happens.

### 5. Apply only what they approved

After an explicit yes, run `fort --only <approved-ids> --fix` (or add `--yes` once
they've confirmed, to skip the second prompt). Fix exactly the checks they agreed to —
never expand the scope. Some fixes need `sudo` and will prompt for a password; that's
the user's machine asking, which is correct.

### 6. Verify and, if useful, document

Re-run `fort --only <ids> --json` to confirm each change took effect. If the user
needs proof for an audit or wants a record, offer `fort --report` and tell them where
the HTML lands (`fort-report.html`, printable to PDF).

## The 16 checks

Use these exact IDs with `--only`. Never invent IDs; if unsure, run `fort --json` and
read them back. Full explainers are in `reference/checks.md`.

- **Core protection:** `passwordmgr`, `filevault`, `screenlock`, `antivirus`
- **System hardening:** `firewall`, `gatekeeper`, `sip`, `ssh`
- **Account controls:** `localadmin`, `guestaccount`, `autologin`, `sudo_touchid`
- **Exposure reduction:** `sharing`, `airdrop`
- **Patching:** `osupdates`, `osversion`

Fixable today: `screenlock`, `firewall`, `gatekeeper`, `ssh`, `guestaccount`,
`autologin`, `sudo_touchid`, `airdrop`, `osupdates`. The rest are audit-only because
the safe action depends on the person (e.g. you can't auto-strip someone's admin
rights or install an antivirus for them) — explain and guide instead.

## Compliance evidence (only when asked)

The audit comes first; framework mapping is a bonus. If the user is preparing for
SOC 2, ISO 27001, NIST CSF, or CIS, each finding already carries its control mappings
in the JSON (`frameworks`). Reference them when relevant ("this is SOC 2 CC6.6"),
and offer `fort --report` for an auditor-ready HTML record. Don't lead with this for
someone who just wants to know if their laptop is locked down.

## Examples

- "Is my Mac secure?" → `fort --json`, give the score and verdict, walk the failures
  by real-world risk, offer to fix the safe ones.
- "Turn on my firewall." → `fort --only firewall --dry-run`, show the command,
  confirm, then `fort --only firewall --fix`, then re-audit to confirm it passes.
- "Fix what's safe." → audit, propose the specific fixable checks (by ID) with a
  one-line reason each, get a yes, then `fort --only <those-ids> --fix`. Leave the
  judgment-call items (admin rights, antivirus) for the user to decide.
- "What breaks if I disable SSH?" → read `reference/checks.md` for `ssh`, explain the
  trade-off, don't change anything unless they ask.
- "I need endpoint evidence for our SOC 2 audit." → `fort --report`, return the path,
  note it prints to PDF, and summarize which controls are covered.
