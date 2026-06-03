//go:build darwin

package checks

import (
	"strings"
)

// FileVaultCheck verifies that full disk encryption is enabled.
type FileVaultCheck struct{}

func (c *FileVaultCheck) ID() string   { return "filevault" }
func (c *FileVaultCheck) Name() string { return "Disk encryption" }

func (c *FileVaultCheck) Run() Result {
	raw, transcript, err := evidenceCmd("fdesetup", "status")
	current, status := "unknown", StatusWarn
	if err == nil {
		switch {
		case strings.Contains(raw, "FileVault is On"):
			current, status = "on", StatusPass
		case strings.Contains(raw, "FileVault is Off"):
			current, status = "off", StatusFail
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: status, Current: current, Expected: "on",
		Evidence: transcript,
	}
}

func (c *FileVaultCheck) Fixable() bool         { return false }
func (c *FileVaultCheck) Fix() error             { return nil }
func (c *FileVaultCheck) FixDescription() string { return "" }
