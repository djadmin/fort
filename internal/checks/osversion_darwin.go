//go:build darwin

package checks

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// OSVersionCheck reports the current macOS version and checks whether
// any software updates are pending according to the last cached update check.
// This complements osupdates (which checks if auto-update is enabled) by
// verifying that updates have actually been applied (CC6.8).
type OSVersionCheck struct{}

func (c *OSVersionCheck) ID() string   { return "osversion" }
func (c *OSVersionCheck) Name() string { return "OS patch status" }

func (c *OSVersionCheck) Run() Result {
	ver := currentOSVersion()

	// Check pending update count from cached SoftwareUpdate prefs (no network call).
	pending := pendingUpdateCount()
	if pending > 0 {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status:   StatusFail,
			Current:  fmt.Sprintf("%s (%d update(s) pending)", ver, pending),
			Expected: "up to date",
			Fixable:  false,
		}
	}

	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusPass, Current: ver, Expected: "up to date",
	}
}

func currentOSVersion() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func pendingUpdateCount() int {
	out, err := exec.Command(
		"defaults", "read",
		"/Library/Preferences/com.apple.SoftwareUpdate",
		"PendingUpdateCount",
	).Output()
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return n
}

func (c *OSVersionCheck) Fixable() bool        { return false }
func (c *OSVersionCheck) Fix() error            { return nil }
func (c *OSVersionCheck) FixDescription() string { return "" }
