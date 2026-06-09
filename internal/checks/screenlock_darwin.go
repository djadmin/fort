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
	askRaw, askT, err1 := evidenceCmd("defaults", "read", "com.apple.screensaver", "askForPassword")
	delayRaw, delayT, err2 := evidenceCmd("defaults", "read", "com.apple.screensaver", "askForPasswordDelay")
	ev := joinTranscripts(askT, delayT)

	// On macOS Ventura and later this legacy preference key is usually unreadable (the
	// setting moved into a protected domain). If we can't read askForPassword at all, we
	// must NOT claim the lock is off — that's a false alarm on a perfectly fine Mac.
	// Report "unknown" and let the user verify in System Settings.
	if err1 != nil {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusWarn, Current: "unknown (verify in System Settings)", Expected: "immediate",
			Fixable: true, Evidence: ev,
		}
	}

	askEnabled := strings.TrimSpace(askRaw) == "1"
	delayImmediate := err2 == nil && strings.TrimSpace(delayRaw) == "0"

	if askEnabled && delayImmediate {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "immediate", Expected: "immediate",
			Fixable: true, Evidence: ev,
		}
	}
	current := "off"
	if askEnabled && !delayImmediate {
		current = "delayed"
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: current, Expected: "immediate",
		Fixable: true, Evidence: ev,
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
