//go:build darwin

package checks

import (
	"fmt"
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
	verRaw, verT, _ := evidenceCmd("sw_vers", "-productVersion")
	ver := strings.TrimSpace(verRaw)
	if ver == "" {
		ver = "unknown"
	}

	pendingRaw, pendingT, _ := evidenceCmd(
		"defaults", "read",
		"/Library/Preferences/com.apple.SoftwareUpdate",
		"PendingUpdateCount",
	)
	pending, _ := strconv.Atoi(strings.TrimSpace(pendingRaw))
	ev := joinTranscripts(verT, pendingT)

	if pending > 0 {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status:   StatusFail,
			Current:  fmt.Sprintf("%s (%d update(s) pending)", ver, pending),
			Expected: "up to date",
			Evidence: ev,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusPass, Current: ver, Expected: "up to date",
		Evidence: ev,
	}
}

func (c *OSVersionCheck) Fixable() bool         { return false }
func (c *OSVersionCheck) Fix() error             { return nil }
func (c *OSVersionCheck) FixDescription() string { return "" }
