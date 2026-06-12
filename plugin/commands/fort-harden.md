---
description: Audit with fort, then fix the safe issues one at a time with your approval
argument-hint: "[check ids to focus on, optional]"
tools:
  allowed: ["Bash", "Read"]
---

Audit this Mac and guide the user through fixing what's safe — **with explicit approval
before any change**. This command can modify system settings, so move carefully.

1. Audit first: `fort --json` (or `fort --only $ARGUMENTS --json` if IDs were given).
2. Identify the **fixable** failing checks (`fixable: true`, `status: fail`). Present
   them as a short list, each with a one-line reason it matters and what fixing it
   would affect. Leave judgment-call items (admin rights, antivirus, OS upgrade) out of
   the fix list — explain those and let the user decide.
3. For each fix the user is interested in, run `fort --only <id> --dry-run` and show the
   exact command fort would run. Let them see it before it happens.
4. Get an explicit yes for the specific checks. Then apply exactly those:
   `fort --only <approved-ids> --fix` (sudo may prompt for a password — that's expected).
   **Never expand beyond what they approved, and never fix everything implicitly.**
5. Re-run `fort --only <ids> --json` to confirm each change took effect, and report what
   changed.
6. Respect the user's setup: if they say they use SSH, need an admin account, etc., do
   not push to change it. Their workflow wins over a generic checklist.

Use `reference/checks.md` from the fort skill for the trade-offs of each fix.

$ARGUMENTS
