//go:build darwin

package checks

import (
	"strings"
)

// SIPCheck verifies System Integrity Protection is enabled.
// SIP prevents root-level modification of system files and is required for
// Gatekeeper and XProtect to function correctly (CC6.8).
type SIPCheck struct{}

func (c *SIPCheck) ID() string   { return "sip" }
func (c *SIPCheck) Name() string { return "System integrity (SIP)" }

func (c *SIPCheck) Run() Result {
	raw, transcript, err := evidenceCmd("csrutil", "status")
	if err != nil {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusWarn, Current: "unknown", Expected: "enabled",
			Evidence: transcript,
		}
	}
	if strings.Contains(raw, "enabled") {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "enabled", Expected: "enabled",
			Evidence: transcript,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: "disabled", Expected: "enabled",
		Evidence: transcript,
	}
}

// SIP can only be re-enabled by booting into Recovery Mode.
func (c *SIPCheck) Fixable() bool         { return false }
func (c *SIPCheck) Fix() error             { return nil }
func (c *SIPCheck) FixDescription() string { return "" }
