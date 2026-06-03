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
	// launchctl print doesn't need sudo and works on macOS 12+.
	// Non-zero exit means the service is not loaded → SSH is off → good.
	_, transcript, err := evidenceCmd("launchctl", "print", "system/com.openssh.sshd")
	if err != nil {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "off", Expected: "off",
			Fixable: true, Evidence: transcript,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: "on", Expected: "off",
		Fixable: true, Evidence: transcript,
	}
}

func (c *SSHCheck) Fixable() bool { return true }

func (c *SSHCheck) Fix() error {
	out, err := exec.Command("sudo", "systemsetup", "-setremotelogin", "off").CombinedOutput()
	if err != nil && strings.Contains(string(out), "not recognized") {
		return exec.Command("sudo", "launchctl", "disable", "system/com.openssh.sshd").Run()
	}
	return err
}

func (c *SSHCheck) FixDescription() string {
	return "sudo systemsetup -setremotelogin off"
}
