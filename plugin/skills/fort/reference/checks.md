# fort checks — what each one means and why it matters

Reference for explaining findings and fixes in plain language. `fort` is the source of
truth for live status; read the `evidence` field in `fort --json` for the actual
command and output behind any result. Get the exact fix command for any check with
`fort --only <id> --dry-run`.

For each check: what it is, why it matters, what to watch out for if you change it.

---

## Core protection

### `passwordmgr` — Password manager *(audit-only)*
**What:** Whether a known password manager (1Password, Bitwarden, etc.) is installed.
**Why:** Reused and weak passwords are the most common way accounts get taken over. A
manager makes unique strong passwords the path of least resistance.
**Changing it:** Can't be auto-installed for the user — installing and adopting a
manager is a personal choice. Guide them to one if it's missing.

### `filevault` — Disk encryption *(audit-only)*
**What:** Whether FileVault full-disk encryption is on.
**Why:** Without it, anyone who gets the physical machine can pull data off the drive.
This is the single most important protection for a laptop that leaves the house.
**Changing it:** Turning FileVault on generates a recovery key — the user must save it
somewhere safe, because losing both the password and the key means the data is gone.
Because of that, fort guides the user through enabling it rather than flipping it
silently.

### `screenlock` — Screen lock *(fixable)*
**What:** Whether the screen locks and requires a password, and how soon after sleep
or screensaver. "immediate" is the target.
**Why:** A delay means a walk-away laptop is open to whoever is nearby for that window.
**Changing it:** Sets the lock to require a password immediately. Harmless; the only
effect is you'll be asked for your password sooner after the screen sleeps.

### `antivirus` — Antivirus / EDR *(audit-only)*
**What:** Whether endpoint protection (CrowdStrike, SentinelOne, etc.) is present.
**Why:** Defense in depth against malware and a common SOC 2 / customer-security-review
requirement for company machines.
**Changing it:** Can't be installed automatically — usually a managed/company decision.
Note that built-in macOS protections (Gatekeeper, XProtect, SIP) still apply even
without a third-party agent.

---

## System hardening

### `firewall` — Application firewall *(fixable)*
**What:** macOS's built-in application firewall (on/off).
**Why:** Blocks unsolicited inbound connections to apps and services on the machine,
shrinking what's reachable on untrusted networks (cafés, airports, hotels).
**Changing it:** `sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate on`.
Low risk. It controls *inbound* connections only; normal outbound use (browsing, etc.)
is unaffected. If the user runs a local server others connect to, macOS may prompt to
allow it the first time.

### `gatekeeper` — Gatekeeper *(fixable)*
**What:** macOS's check that apps are signed and notarized before they run.
**Why:** First line of defense against running tampered or malicious downloaded apps.
**Changing it:** Re-enables Gatekeeper. Watch-out: if the user routinely runs unsigned
or homebrew-built apps, they'll now get a prompt and need to allow those explicitly.
Worth a heads-up for developers.

### `sip` — System Integrity Protection *(audit-only)*
**What:** Kernel-level protection that stops even root from modifying protected system
files.
**Why:** A major barrier to malware and accidental system corruption persisting.
**Changing it:** Re-enabling SIP requires booting into Recovery Mode — it can't be done
from a normal session, so fort reports it rather than fixing it. If it's off, the user
likely disabled it deliberately (some dev tooling needs it off); confirm before
advising they re-enable.

### `ssh` — Remote login (SSH) *(fixable)*
**What:** Whether the SSH remote-login service is accepting connections.
**Why:** An open SSH service is a remotely reachable way into the machine and a brute-
force target if exposed.
**Changing it:** The fix turns remote login off. **Ask first** — if the user SSHes into
this Mac (headless box, remote dev, file sync), disabling it cuts that off. If they
never use it, turning it off is a clear win.

---

## Account controls

### `localadmin` — Local admin rights *(audit-only)*
**What:** Whether the day-to-day account has admin privileges.
**Why:** Running as a standard user limits the blast radius if the account is
compromised or tricked into running something — malware can't silently make
system-wide changes.
**Changing it:** Can't be auto-fixed: stripping admin from the only admin account can
lock the user out of installing software and system changes. The right move is usually
a separate admin account plus a standard daily account — explain that rather than
flipping it.

### `guestaccount` — Guest account *(fixable)*
**What:** Whether the macOS Guest user is enabled.
**Why:** An enabled guest account is an unauthenticated local entry point.
**Changing it:** Disables Guest login. Safe for almost everyone; only matters if someone
deliberately relies on a guest session.

### `autologin` — Automatic login *(fixable)*
**What:** Whether the Mac logs into an account automatically at boot, with no password.
**Why:** Defeats the login password entirely — power on equals full access.
**Changing it:** Disables auto-login so a password is required at boot. The only effect
is the user types their password when the machine starts. Note: FileVault already
forces this, so the two interact.

### `sudo_touchid` — Touch ID for sudo *(fixable)*
**What:** Whether Touch ID can authenticate `sudo` in the terminal.
**Why:** A usability and security win — biometric prompts beat typed passwords that can
be shoulder-surfed or muscle-memoried into the wrong window.
**Changing it:** Enables Touch ID for sudo. Low risk; password sudo still works as a
fallback. Only relevant on Macs with a Touch ID sensor.

---

## Exposure reduction

### `sharing` — Sharing services *(audit-only)*
**What:** Whether macOS sharing services (file, screen, printer, media sharing, etc.)
are off.
**Why:** Each enabled service is another listening door on the network.
**Changing it:** fort audits these rather than batch-disabling, because some are
intentional (a user actively sharing their screen or files). Review which are on and
turn off the ones that aren't needed, individually.

### `airdrop` — AirDrop *(fixable)*
**What:** AirDrop receiving mode — Off, Contacts Only, or Everyone.
**Why:** "Everyone" lets nearby strangers send the machine files and prompts; "Contacts
Only" or "Off" closes that.
**Changing it:** Tightens AirDrop to a safer mode. Effect: the user may need to widen it
again temporarily when receiving a file from someone not in their contacts.

---

## Patching

### `osupdates` — Automatic OS updates *(fixable)*
**What:** Whether macOS automatically checks for and installs security updates.
**Why:** Most real-world compromises use already-patched bugs. Auto-updates close that
gap without the user having to remember.
**Changing it:** Enables automatic update checks and security-response installs. Low
risk; the main trade-off is updates installing on Apple's schedule rather than the
user's. Reassure update-shy users that security responses are small and reversible.

### `osversion` — OS patch status *(audit-only)*
**What:** Whether the installed macOS version is current / still receiving security
fixes.
**Why:** An OS past its support window stops getting security patches entirely.
**Changing it:** Can't be auto-fixed — a major OS upgrade is the user's call and can
take time and disk space. If they're behind, point them to System Settings → Software
Update and note whether their hardware supports the current release.
