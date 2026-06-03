//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

const sockfilterfw = "/usr/libexec/ApplicationFirewall/socketfilterfw"

// FirewallCheck verifies the macOS application firewall is enabled.
type FirewallCheck struct{}

func (c *FirewallCheck) ID() string   { return "firewall" }
func (c *FirewallCheck) Name() string { return "Application firewall" }

func (c *FirewallCheck) Run() Result {
	raw, transcript, err := evidenceCmd(sockfilterfw, "--getglobalstate")
	current, status := "unknown", StatusWarn
	if err == nil {
		s := strings.ToLower(raw)
		switch {
		case strings.Contains(s, "enabled") || strings.Contains(s, "state = 1"):
			current, status = "on", StatusPass
		case strings.Contains(s, "disabled") || strings.Contains(s, "state = 0"):
			current, status = "off", StatusFail
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: status, Current: current, Expected: "on",
		Fixable: true, Evidence: transcript,
	}
}

func (c *FirewallCheck) Fixable() bool { return true }

func (c *FirewallCheck) Fix() error {
	return exec.Command("sudo", sockfilterfw, "--setglobalstate", "on").Run()
}

func (c *FirewallCheck) FixDescription() string {
	return "sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate on"
}
