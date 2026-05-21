//go:build darwin

package checks

import (
	"os"
	"os/exec"
)

var avApps = []struct {
	app  string
	name string
}{
	{"CrowdStrike Falcon.app", "CrowdStrike Falcon"},
	{"Falcon.app", "CrowdStrike Falcon"},
	{"Malwarebytes.app", "Malwarebytes"},
	{"Sophos Anti-Virus.app", "Sophos"},
	{"Sophos Endpoint.app", "Sophos Endpoint"},
	{"ESET Endpoint Security.app", "ESET"},
	{"Bitdefender.app", "Bitdefender"},
	{"Bitdefender Antivirus for Mac.app", "Bitdefender"},
	{"Trend Micro Security.app", "Trend Micro"},
	{"Avast Security.app", "Avast"},
	{"AVG AntiVirus.app", "AVG"},
	{"Carbon Black.app", "VMware Carbon Black"},
	{"SentinelOne.app", "SentinelOne"},
	{"Cylance Desktop.app", "Cylance"},
	{"Symantec Endpoint Protection.app", "Symantec"},
	{"Norton 360.app", "Norton"},
	{"Kaspersky Internet Security.app", "Kaspersky"},
}

var avProcesses = []struct {
	proc string
	name string
}{
	{"falconctl", "CrowdStrike Falcon"},
	{"SentinelAgent", "SentinelOne"},
	{"cbsecurity", "Carbon Black"},
	{"MalwarebytesAgent", "Malwarebytes"},
	{"SophosScanD", "Sophos"},
}

// AntivirusCheck detects whether a third-party AV or EDR agent is installed.
type AntivirusCheck struct{}

func (c *AntivirusCheck) ID() string   { return "antivirus" }
func (c *AntivirusCheck) Name() string { return "Antivirus / EDR" }

func (c *AntivirusCheck) Run() Result {
	for _, av := range avApps {
		if _, err := os.Stat("/Applications/" + av.app); err == nil {
			return Result{
				ID: c.ID(), Name: c.Name(),
				Status: StatusPass, Current: av.name, Expected: "installed",
			}
		}
	}
	for _, av := range avProcesses {
		if err := exec.Command("pgrep", "-x", av.proc).Run(); err == nil {
			return Result{
				ID: c.ID(), Name: c.Name(),
				Status: StatusPass, Current: av.name, Expected: "installed",
			}
		}
	}
	// XProtect is always present on macOS — minimal built-in protection only
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusWarn, Current: "XProtect only (built-in)", Expected: "third-party AV/EDR",
	}
}

func (c *AntivirusCheck) Fixable() bool        { return false }
func (c *AntivirusCheck) Fix() error            { return nil }
func (c *AntivirusCheck) FixDescription() string { return "" }
