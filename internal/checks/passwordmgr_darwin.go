//go:build darwin

package checks

import (
	"os"
	"path/filepath"
)

var knownPasswordManagers = []struct {
	app  string
	name string
}{
	{"1Password.app", "1Password"},
	{"1Password 7 - Password Manager.app", "1Password 7"},
	{"Bitwarden.app", "Bitwarden"},
	{"LastPass.app", "LastPass"},
	{"Dashlane.app", "Dashlane"},
	{"Dashlane - Password Manager.app", "Dashlane"},
	{"Keeper Password Manager.app", "Keeper"},
	{"Enpass.app", "Enpass"},
	{"NordPass.app", "NordPass"},
	{"RoboForm.app", "RoboForm"},
	{"Strongbox.app", "Strongbox"},
	{"KeePassXC.app", "KeePassXC"},
	{"Secrets.app", "Secrets"},
}

// PasswordMgrCheck detects whether a dedicated password manager is installed.
type PasswordMgrCheck struct{}

func (c *PasswordMgrCheck) ID() string   { return "passwordmgr" }
func (c *PasswordMgrCheck) Name() string { return "Password manager" }

func (c *PasswordMgrCheck) Run() Result {
	searchPaths := []string{
		"/Applications",
		filepath.Join(os.Getenv("HOME"), "Applications"),
	}
	for _, pm := range knownPasswordManagers {
		for _, dir := range searchPaths {
			path := filepath.Join(dir, pm.app)
			if _, err := os.Stat(path); err == nil {
				_, transcript, _ := evidenceCmd("ls", "-d", path)
				return Result{
					ID: c.ID(), Name: c.Name(),
					Status: StatusPass, Current: pm.name, Expected: "installed",
					Evidence: transcript,
				}
			}
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: "none found", Expected: "installed",
		Evidence: "# Checked /Applications and ~/Applications — no known password manager found",
	}
}

func (c *PasswordMgrCheck) Fixable() bool         { return false }
func (c *PasswordMgrCheck) Fix() error             { return nil }
func (c *PasswordMgrCheck) FixDescription() string { return "" }
