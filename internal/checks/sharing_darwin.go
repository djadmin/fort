//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// SharingCheck verifies all inbound sharing services are off:
// file sharing (SMB), screen sharing, remote management (ARD), internet sharing.
// Each is an attack surface that should be off on employee laptops (CC6.6).
type SharingCheck struct{}

func (c *SharingCheck) ID() string   { return "sharing" }
func (c *SharingCheck) Name() string { return "Sharing services" }

func (c *SharingCheck) Run() Result {
	active := c.activeServices()
	if len(active) == 0 {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "all off", Expected: "all off",
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail,
		Current:  strings.Join(active, ", "),
		Expected: "all off",
		Fixable:  false,
	}
}

func (c *SharingCheck) activeServices() []string {
	var on []string
	if launchdRunning("com.apple.smbd") {
		on = append(on, "file sharing")
	}
	if launchdRunning("com.apple.screensharing") {
		on = append(on, "screen sharing")
	}
	if launchdRunning("com.apple.RemoteDesktopAgent") {
		on = append(on, "remote management")
	}
	if launchdRunning("com.apple.InternetSharing") {
		on = append(on, "internet sharing")
	}
	return on
}

func launchdRunning(service string) bool {
	return exec.Command("launchctl", "list", service).Run() == nil
}

func (c *SharingCheck) Fixable() bool        { return false }
func (c *SharingCheck) Fix() error            { return nil }
func (c *SharingCheck) FixDescription() string { return "" }
