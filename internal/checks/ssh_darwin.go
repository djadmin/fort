//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// SSHCheck verifies Remote Login (inbound SSH) is disabled on the endpoint.
type SSHCheck struct{}

func (c *SSHCheck) ID() string   { return "ssh" }
func (c *SSHCheck) Name() string { return "Remote login (SSH)" }

func (c *SSHCheck) Run() Result {
	// launchctl print doesn't need sudo and works on macOS 12+
	err := exec.Command("launchctl", "print", "system/com.openssh.sshd").Run()
	if err != nil {
		// service not loaded → SSH is off → good
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "off", Expected: "off", Fixable: true,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: "on", Expected: "off", Fixable: true,
	}
}

func (c *SSHCheck) Fixable() bool { return true }

func (c *SSHCheck) Fix() error {
	// systemsetup still works on macOS 12+ (deprecated warning is non-fatal)
	out, err := exec.Command("sudo", "systemsetup", "-setremotelogin", "off").CombinedOutput()
	if err != nil && strings.Contains(string(out), "not recognized") {
		// Fallback for newer macOS
		return exec.Command("sudo", "launchctl", "disable", "system/com.openssh.sshd").Run()
	}
	return err
}

func (c *SSHCheck) FixDescription() string {
	return "sudo systemsetup -setremotelogin off"
}
