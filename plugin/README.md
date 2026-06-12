# fort — Claude Code plugin

Audit and harden your Mac's security from a conversation. This plugin teaches Claude
Code to drive the [`fort`](https://github.com/djadmin/fort) CLI: run a full security
audit, explain every finding in plain language (what it is, why it matters, what
changing it would affect), and fix only what you approve.

It runs the `fort` binary over your shell — no MCP server, no extra daemon. You just
need `fort` installed.

## Install

```bash
# 1. install the fort CLI (the plugin orchestrates this binary)
brew install djadmin/tap/fort

# 2. add this repo as a plugin marketplace, then install the plugin
/plugin marketplace add djadmin/fort
/plugin install fort@fort
```

## Use it

Just ask:

- "Is my Mac secure?" — runs the audit and walks the findings by real-world risk.
- "Fix what's safe." — proposes the safe fixes, previews each, applies only what you OK.
- "Turn on my firewall." — previews the exact command, confirms, applies, re-checks.
- "What breaks if I disable SSH?" — explains the trade-off, changes nothing.
- "I need endpoint evidence for our SOC 2 audit." — generates an HTML report.

Or use the commands directly:

| Command | What it does | Touches your system? |
|---|---|---|
| `/fort-audit` | Audit and explain the findings | No (read-only) |
| `/fort-harden` | Audit, then fix safe issues with your approval | Only what you approve |
| `/fort-report` | Write an auditor-ready HTML evidence report | No (writes a file) |

## Safety

Auditing is always read-only. Fixes change system settings, so the plugin previews the
exact command, waits for your explicit yes, and applies only the specific checks you
approved — never "everything" implicitly. `fort` itself is the source of truth for what
your posture is; the plugin adds the explanation and judgment around it.
