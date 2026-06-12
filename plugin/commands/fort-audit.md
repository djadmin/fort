---
description: Audit this Mac's security posture with fort and explain the findings in plain language
argument-hint: "[check ids, e.g. firewall,filevault]"
tools:
  allowed: ["Bash", "Read"]
---

Run a read-only macOS security audit with `fort` and present it so the user actually
understands where they stand. **Do not change anything in this command** — auditing is
read-only.

1. Confirm fort is installed (`fort --version`); if not, tell the user
   `brew install djadmin/tap/fort` and stop.
2. Run `fort --json`. If the user passed check IDs in `$ARGUMENTS`, scope it:
   `fort --only $ARGUMENTS --json`.
3. Lead with the score and a one-line verdict.
4. Walk the findings in order of real-world risk — what exposes the machine first,
   lower-risk items after. For each failing or warning check, say what it is, where it
   stands (`current` vs `expected`, and the `evidence` fort captured), and why it
   matters. Note which are auto-fixable.
5. Briefly note what passed.
6. Offer next steps: preview a fix (`fort --only <id> --dry-run`) or, if they want it,
   an evidence report (`fort --report`). Don't apply anything yet — fixes happen only
   after the user reviews a preview and explicitly says yes.

Use `reference/checks.md` from the fort skill for the "why it matters" detail.

$ARGUMENTS
