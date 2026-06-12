---
description: Generate an auditor-ready HTML evidence report of this Mac's security posture
tools:
  allowed: ["Bash", "Read"]
---

Generate a security evidence report with `fort` for an auditor or a personal record.
This is read-only — it writes a file, it doesn't change any system settings.

1. Confirm fort is installed (`fort --version`); if not, tell the user
   `brew install djadmin/tap/fort` and stop.
2. Run `fort --report`. It writes `fort-report.html` in the current directory.
3. Tell the user the path and that it's a self-contained HTML file they can open in a
   browser or print to PDF for an auditor.
4. Briefly summarize what the report covers: the current posture and the SOC 2 / ISO
   27001 / NIST CSF / CIS v8 control mappings for each check. Lead with the security
   picture; the framework mappings are there for whoever needs the audit trail.

$ARGUMENTS
