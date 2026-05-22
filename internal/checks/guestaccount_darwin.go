//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// GuestAccountCheck verifies the guest account is disabled.
// A guest account allows physical access bypassing all logical controls (CC6.1).
type GuestAccountCheck struct{}

func (c *GuestAccountCheck) ID() string   { return "guestaccount" }
func (c *GuestAccountCheck) Name() string { return "Guest account" }

func (c *GuestAccountCheck) Run() Result {
	out, err := exec.Command(
		"defaults", "read",
		"/Library/Preferences/com.apple.loginwindow",
		"GuestEnabled",
	).Output()

	// Key missing or error → guest account is disabled (default)
	if err != nil || strings.TrimSpace(string(out)) == "0" {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "disabled", Expected: "disabled",
			Fixable: true,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: "enabled", Expected: "disabled",
		Fixable: true,
	}
}

func (c *GuestAccountCheck) Fixable() bool { return true }

func (c *GuestAccountCheck) Fix() error {
	return exec.Command(
		"sudo", "defaults", "write",
		"/Library/Preferences/com.apple.loginwindow",
		"GuestEnabled", "-bool", "false",
	).Run()
}

func (c *GuestAccountCheck) FixDescription() string {
	return "sudo defaults write /Library/Preferences/com.apple.loginwindow GuestEnabled -bool false"
}
