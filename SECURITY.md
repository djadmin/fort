# Security Policy

`fort` is a security tool, so we take the security of the tool itself seriously.

## Reporting a vulnerability

**Please do not open a public issue for security vulnerabilities.**

Report privately through GitHub's [private vulnerability reporting](https://github.com/djadmin/fort/security/advisories/new). Click **"Report a vulnerability"** under the repository's **Security** tab. This creates a private advisory only you and the maintainers can see.

If you can't use GitHub advisories, email the maintainer (see the GitHub profile at [github.com/djadmin](https://github.com/djadmin)) with the subject line `fort security`.

When reporting, please include:

- the version of `fort` (`fort --version`) and your macOS version
- a description of the issue and its impact
- steps to reproduce, ideally with the exact command and output

## What to expect

- We aim to acknowledge a report within 3 business days.
- We'll confirm the issue, work on a fix, and keep you updated on progress.
- Once a fix ships, we'll credit you in the release notes unless you'd rather stay anonymous.

## Scope

`fort` reads local macOS security state and, with `--fix`, changes documented settings using public Apple interfaces. Reports we're especially interested in:

- a check that reports a **false pass** (says a control is on when it isn't)
- a `--fix` that makes an unintended or unsafe change to the system
- command injection or privilege issues in how checks shell out to system tools
- anything that could cause `fort` to leak data off the machine (it shouldn't, there's no network call in the audit path)

## Supported versions

Security fixes are applied to the latest released version. Please upgrade (`brew upgrade djadmin/tap/fort`) before reporting.
