//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// AirDropCheck verifies AirDrop is not set to Everyone.
// AirDrop set to Everyone allows any nearby device to send files (CC6.6).
type AirDropCheck struct{}

func (c *AirDropCheck) ID() string   { return "airdrop" }
func (c *AirDropCheck) Name() string { return "AirDrop" }

func (c *AirDropCheck) Run() Result {
	out, err := exec.Command(
		"defaults", "read",
		"com.apple.sharingd",
		"DiscoverableMode",
	).Output()

	if err != nil {
		// Key missing on newer macOS — try reading via system-level domain
		out2, err2 := exec.Command(
			"defaults", "read",
			"/Library/Preferences/com.apple.sharingd",
			"DiscoverableMode",
		).Output()
		if err2 != nil {
			return Result{
				ID: c.ID(), Name: c.Name(),
				Status: StatusWarn, Current: "unknown", Expected: "Contacts Only or Off",
			}
		}
		out = out2
	}

	mode := strings.TrimSpace(string(out))
	switch mode {
	case "Off", "Contacts Only":
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: mode, Expected: "Contacts Only or Off",
			Fixable: true,
		}
	case "Everyone":
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusFail, Current: mode, Expected: "Contacts Only or Off",
			Fixable: true,
		}
	default:
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusWarn, Current: mode, Expected: "Contacts Only or Off",
		}
	}
}

func (c *AirDropCheck) Fixable() bool { return true }

func (c *AirDropCheck) Fix() error {
	if err := exec.Command(
		"defaults", "write", "com.apple.sharingd",
		"DiscoverableMode", "-string", "Contacts Only",
	).Run(); err != nil {
		return err
	}
	// Restart sharingd to apply immediately
	_ = exec.Command("killall", "sharingd").Run()
	return nil
}

func (c *AirDropCheck) FixDescription() string {
	return `defaults write com.apple.sharingd DiscoverableMode -string "Contacts Only" && killall sharingd`
}
