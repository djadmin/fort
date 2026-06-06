//go:build darwin

package checks

// All returns all security checks for macOS, grouped by category.
func All() []Check {
	return []Check{
		// ── Core security (4) ──────────────────────────────────────────
		&PasswordMgrCheck{},  // CC6.1 — credential hygiene
		&FileVaultCheck{},    // CC6.7 — encryption at rest
		&ScreenLockCheck{},   // CC6.1 — physical access control
		&AntivirusCheck{},    // CC6.8 — malware prevention

		// ── System hardening (4) ───────────────────────────────────────
		&FirewallCheck{},     // CC6.6 — boundary protection
		&GatekeeperCheck{},   // CC6.8 — software allowlisting
		&SIPCheck{},          // CC6.8 — system integrity
		&SSHCheck{},          // CC6.6 — remote access control

		// ── Access controls (4) ────────────────────────────────────────
		&LocalAdminCheck{},    // CC6.3 — least privilege
		&GuestAccountCheck{},  // CC6.1 — unauthorized access prevention
		&AutoLoginCheck{},     // CC6.1 — physical access bypass prevention
		&SudoTouchIDCheck{},   // CC6.1 — privileged command authentication

		// ── Exposure reduction (2) ─────────────────────────────────────
		&SharingCheck{},      // CC6.6 — attack surface reduction
		&AirDropCheck{},      // CC6.6 — data exfiltration prevention

		// ── Patching (2) ──────────────────────────────────────────────
		&OSUpdatesCheck{},    // CC6.8 — automatic update policy
		&OSVersionCheck{},    // CC6.8 — patch application verification
	}
}
