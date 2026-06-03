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
	raw, transcript, err := evidenceCmd(
		"defaults", "read",
		"com.apple.sharingd",
		"DiscoverableMode",
	)
	if err != nil {
		// Key missing on newer macOS — try the system-level domain.
		raw, transcript, err = evidenceCmd(
			"defaults", "read",
			"/Library/Preferences/com.apple.sharingd",
			"DiscoverableMode",
		)
		if err != nil {
			return Result{
				ID: c.ID(), Name: c.Name(),
				Status: StatusWarn, Current: "unknown", Expected: "Contacts Only or Off",
				Fixable: true, Evidence: transcript,
			}
		}
	}

	mode := strings.TrimSpace(raw)
	switch mode {
	case "Off", "Contacts Only":
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: mode, Expected: "Contacts Only or Off",
			Fixable: true, Evidence: transcript,
		}
	case "Everyone":
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusFail, Current: mode, Expected: "Contacts Only or Off",
			Fixable: true, Evidence: transcript,
		}
	default:
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusWarn, Current: mode, Expected: "Contacts Only or Off",
			Fixable: true, Evidence: transcript,
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
	_ = exec.Command("killall", "sharingd").Run()
	return nil
}

func (c *AirDropCheck) FixDescription() string {
	return `defaults write com.apple.sharingd DiscoverableMode -string "Contacts Only" && killall sharingd`
}
