//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// OSUpdatesCheck verifies automatic update checks are enabled.
type OSUpdatesCheck struct{}

func (c *OSUpdatesCheck) ID() string   { return "osupdates" }
func (c *OSUpdatesCheck) Name() string { return "Automatic OS updates" }

func (c *OSUpdatesCheck) Run() Result {
	out, err := exec.Command(
		"defaults", "read",
		"/Library/Preferences/com.apple.SoftwareUpdate",
		"AutomaticCheckEnabled",
	).Output()

	enabled := err == nil && strings.TrimSpace(string(out)) == "1"
	if enabled {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "on", Expected: "on", Fixable: true,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: "off", Expected: "on", Fixable: true,
	}
}

func (c *OSUpdatesCheck) Fixable() bool { return true }

func (c *OSUpdatesCheck) Fix() error {
	return exec.Command(
		"sudo", "defaults", "write",
		"/Library/Preferences/com.apple.SoftwareUpdate",
		"AutomaticCheckEnabled", "-bool", "true",
	).Run()
}

func (c *OSUpdatesCheck) FixDescription() string {
	return "sudo defaults write /Library/Preferences/com.apple.SoftwareUpdate AutomaticCheckEnabled -bool true"
}
