//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// SIPCheck verifies System Integrity Protection is enabled.
// SIP prevents root-level modification of system files and is required for
// Gatekeeper and XProtect to function correctly (CC6.8).
type SIPCheck struct{}

func (c *SIPCheck) ID() string   { return "sip" }
func (c *SIPCheck) Name() string { return "System integrity (SIP)" }

func (c *SIPCheck) Run() Result {
	out, err := exec.Command("csrutil", "status").Output()
	if err != nil {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusWarn, Current: "unknown", Expected: "enabled",
		}
	}
	s := strings.TrimSpace(string(out))
	if strings.Contains(s, "enabled") {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "enabled", Expected: "enabled",
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: "disabled", Expected: "enabled",
		Fixable: false,
	}
}

// SIP can only be re-enabled by booting into Recovery Mode.
func (c *SIPCheck) Fixable() bool        { return false }
func (c *SIPCheck) Fix() error            { return nil }
func (c *SIPCheck) FixDescription() string { return "" }
