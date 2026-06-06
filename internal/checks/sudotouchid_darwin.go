//go:build darwin

package checks

import (
	"os/exec"
	"strings"
)

// SudoTouchIDCheck verifies that PAM is configured to accept TouchID for sudo.
// macOS 14+ (Sonoma): the fix writes /etc/pam.d/sudo_local, which is designed
// to be user-managed and survives OS updates. macOS 12–13 uses /etc/pam.d/sudo
// (overwritten on system updates, so the fix is guidance only there).
type SudoTouchIDCheck struct{}

func (c *SudoTouchIDCheck) ID() string   { return "sudo_touchid" }
func (c *SudoTouchIDCheck) Name() string { return "TouchID for sudo" }

func (c *SudoTouchIDCheck) Run() Result {
	// macOS 14+ (Sonoma): /etc/pam.d/sudo includes sudo_local; user creates
	// sudo_local to enable TouchID. Apple ships sudo_local.template for guidance.
	// macOS 12–13: pam_tid.so can appear directly in /etc/pam.d/sudo.
	localOut, localT, _ := evidenceCmd("cat", "/etc/pam.d/sudo_local")
	sudoOut, sudoT, _ := evidenceCmd("cat", "/etc/pam.d/sudo")
	tmplOut, tmplT, _ := evidenceCmd("cat", "/etc/pam.d/sudo_local.template")
	ev := joinTranscripts(localT, sudoT, tmplT)
	_ = tmplOut

	if hasTouchIDLine(localOut) || hasTouchIDLine(sudoOut) {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "enabled", Expected: "enabled",
			Fixable: true, Evidence: ev,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status: StatusFail, Current: "disabled", Expected: "enabled",
		Fixable: true, Evidence: ev,
	}
}

// hasTouchIDLine returns true if the PAM content has an uncommented pam_tid.so line.
func hasTouchIDLine(content string) bool {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "pam_tid.so") {
			return true
		}
	}
	return false
}

func (c *SudoTouchIDCheck) Fixable() bool { return true }

func (c *SudoTouchIDCheck) Fix() error {
	// Write /etc/pam.d/sudo_local — the macOS 14+ user-managed PAM include.
	// "sufficient" means TouchID is accepted; password still works as a fallback.
	cmd := exec.Command("sudo", "tee", "/etc/pam.d/sudo_local")
	cmd.Stdin = strings.NewReader("auth       sufficient     pam_tid.so\n")
	_, err := cmd.Output()
	return err
}

func (c *SudoTouchIDCheck) FixDescription() string {
	return `echo "auth       sufficient     pam_tid.so" | sudo tee /etc/pam.d/sudo_local`
}
