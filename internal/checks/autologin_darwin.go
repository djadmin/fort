//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// AutoLoginCheck verifies automatic login is disabled.
// Auto-login means physical access to the device bypasses all authentication (CC6.1).
type AutoLoginCheck struct{}

func (c *AutoLoginCheck) ID() string   { return "autologin" }
func (c *AutoLoginCheck) Name() string { return "Automatic login" }

func (c *AutoLoginCheck) Run() Result {
	raw, transcript, err := evidenceCmd(
		"defaults", "read",
		"/Library/Preferences/com.apple.loginwindow",
		"autoLoginUser",
	)
	// Error = key doesn't exist = auto-login is disabled
	if err != nil || strings.TrimSpace(raw) == "" {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "disabled", Expected: "disabled",
			Fixable: true, Evidence: transcript,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status:   StatusFail,
		Current:  "enabled for " + strings.TrimSpace(raw),
		Expected: "disabled",
		Fixable:  true, Evidence: transcript,
	}
}

func (c *AutoLoginCheck) Fixable() bool { return true }

func (c *AutoLoginCheck) Fix() error {
	return exec.Command(
		"sudo", "defaults", "delete",
		"/Library/Preferences/com.apple.loginwindow",
		"autoLoginUser",
	).Run()
}

func (c *AutoLoginCheck) FixDescription() string {
	return "sudo defaults delete /Library/Preferences/com.apple.loginwindow autoLoginUser"
}
