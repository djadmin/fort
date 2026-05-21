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
	out, err := exec.Command("spctl", "--status").Output()
	current, status := "unknown", StatusWarn
	if err == nil {
		s := strings.TrimSpace(string(out))
		switch {
		case strings.Contains(s, "enabled"):
			current, status = "enabled", StatusPass
		case strings.Contains(s, "disabled"):
			current, status = "disabled", StatusFail
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: status, Current: current, Expected: "enabled", Fixable: true,
	}
}

func (c *GatekeeperCheck) Fixable() bool { return true }

func (c *GatekeeperCheck) Fix() error {
	return exec.Command("sudo", "spctl", "--master-enable").Run()
}

func (c *GatekeeperCheck) FixDescription() string {
	return "sudo spctl --master-enable"
}
