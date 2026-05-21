//go:build darwin

package checks

// All returns all security checks for macOS, in display order.
func All() []Check {
	return []Check{
		&PasswordMgrCheck{},
		&FileVaultCheck{},
		&ScreenLockCheck{},
		&AntivirusCheck{},
		&FirewallCheck{},
		&GatekeeperCheck{},
		&SSHCheck{},
		&OSUpdatesCheck{},
	}
}
