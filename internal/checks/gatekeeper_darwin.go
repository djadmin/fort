//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// GatekeeperCheck verifies Gatekeeper is enabled (blocks unsigned apps).
type GatekeeperCheck struct{}

func (c *GatekeeperCheck) ID() string   { return "gatekeeper" }
func (c *GatekeeperCheck) Name() string { return "Gatekeeper" }

func (c *GatekeeperCheck) Run() Result {
	// spctl writes to stderr on some macOS versions; CombinedOutput in evidenceCmd catches both.
	raw, transcript, err := evidenceCmd("spctl", "--status")
	current, status := "unknown", StatusWarn
	if err == nil {
		switch {
		case strings.Contains(raw, "enabled"):
			current, status = "enabled", StatusPass
		case strings.Contains(raw, "disabled"):
			current, status = "disabled", StatusFail
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: status, Current: current, Expected: "enabled",
		Fixable: true, Evidence: transcript,
	}
}

func (c *GatekeeperCheck) Fixable() bool { return true }

func (c *GatekeeperCheck) Fix() error {
	return exec.Command("sudo", "spctl", "--master-enable").Run()
}

func (c *GatekeeperCheck) FixDescription() string {
	return "sudo spctl --master-enable"
}
