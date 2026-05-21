//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// ScreenLockCheck verifies the screen locks immediately on sleep/screensaver.
type ScreenLockCheck struct{}

func (c *ScreenLockCheck) ID() string   { return "screenlock" }
func (c *ScreenLockCheck) Name() string { return "Screen lock" }

func (c *ScreenLockCheck) Run() Result {
	askOut, err1 := exec.Command("defaults", "read", "com.apple.screensaver", "askForPassword").Output()
	delayOut, err2 := exec.Command("defaults", "read", "com.apple.screensaver", "askForPasswordDelay").Output()

	askEnabled := err1 == nil && strings.TrimSpace(string(askOut)) == "1"
	delayImmediate := err2 == nil && strings.TrimSpace(string(delayOut)) == "0"

	if askEnabled && delayImmediate {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "immediate", Expected: "immediate", Fixable: true,
		}
	}
	current := "off"
	if askEnabled && !delayImmediate {
		current = "delayed"
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: current, Expected: "immediate", Fixable: true,
	}
}

func (c *ScreenLockCheck) Fixable() bool { return true }

func (c *ScreenLockCheck) Fix() error {
	if err := exec.Command("defaults", "write", "com.apple.screensaver", "askForPassword", "-int", "1").Run(); err != nil {
		return err
	}
	return exec.Command("defaults", "write", "com.apple.screensaver", "askForPasswordDelay", "-int", "0").Run()
}

func (c *ScreenLockCheck) FixDescription() string {
	return "defaults write com.apple.screensaver askForPassword -int 1\n" +
		"     defaults write com.apple.screensaver askForPasswordDelay -int 0"
}
