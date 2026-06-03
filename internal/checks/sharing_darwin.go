//go:build darwin

package checks

import (
	"strings"
)

// SharingCheck verifies all inbound sharing services are off:
// file sharing (SMB), screen sharing, remote management (ARD), internet sharing.
// Each is an attack surface that should be off on employee laptops (CC6.6).
type SharingCheck struct{}

func (c *SharingCheck) ID() string   { return "sharing" }
func (c *SharingCheck) Name() string { return "Sharing services" }

var sharingServices = []struct {
	launchdName string
	label       string
}{
	{"com.apple.smbd", "file sharing"},
	{"com.apple.screensharing", "screen sharing"},
	{"com.apple.RemoteDesktopAgent", "remote management"},
	{"com.apple.InternetSharing", "internet sharing"},
}

func (c *SharingCheck) Run() Result {
	var active, transcripts []string
	for _, svc := range sharingServices {
		_, transcript, err := evidenceCmd("launchctl", "list", svc.launchdName)
		transcripts = append(transcripts, transcript)
		if err == nil {
			active = append(active, svc.label)
		}
	}
	ev := joinTranscripts(transcripts...)
	if len(active) == 0 {
		return Result{
			ID: c.ID(), Name: c.Name(),
			Status: StatusPass, Current: "all off", Expected: "all off",
			Evidence: ev,
		}
	}
	return Result{
		ID: c.ID(), Name: c.Name(),
		Status:   StatusFail,
		Current:  strings.Join(active, ", "),
		Expected: "all off",
		Evidence: ev,
	}
}

func (c *SharingCheck) Fixable() bool         { return false }
func (c *SharingCheck) Fix() error             { return nil }
func (c *SharingCheck) FixDescription() string { return "" }
