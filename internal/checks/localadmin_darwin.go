//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// LocalAdminCheck verifies the current user is not running with local admin rights.
// Running as admin means malware can install system-wide and bypass controls (CC6.3).
type LocalAdminCheck struct{}

func (c *LocalAdminCheck) ID() string   { return "localadmin" }
func (c *LocalAdminCheck) Name() string { return "Local admin rights" }

func (c *LocalAdminCheck) Run() Result {
	out, err := exec.Command("id").Output()
	if err != nil {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusWarn, Current: "unknown", Expected: "standard user",
		}
	}
	if strings.Contains(string(out), "(admin)") {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusFail, Current: "admin", Expected: "standard user",
			Fixable: false,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusPass, Current: "standard user", Expected: "standard user",
	}
}

func (c *LocalAdminCheck) Fixable() bool { return false }
func (c *LocalAdminCheck) Fix() error     { return nil }
func (c *LocalAdminCheck) FixDescription() string { return "" }
