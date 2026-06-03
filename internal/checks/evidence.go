package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

// evidenceCmd runs a command, captures combined stdout+stderr, and returns
// both the raw output and a terminal transcript for embedding in audit reports.
// A non-zero exit is reflected in the transcript but does not panic.
func evidenceCmd(name string, args ...string) (output, transcript string, err error) {
	out, runErr := exec.Command(name, args...).CombinedOutput()
	cmdStr := name
	if len(args) > 0 {
		cmdStr += " " + strings.Join(args, " ")
	}
	raw := strings.TrimSpace(string(out))
	if runErr != nil {
		note := raw
		if note == "" {
			note = "(no output)"
		}
		return raw, fmt.Sprintf("$ %s\n%s", cmdStr, note), runErr
	}
	return raw, fmt.Sprintf("$ %s\n%s", cmdStr, raw), nil
}

// joinTranscripts concatenates transcripts with a blank line between each.
func joinTranscripts(parts ...string) string {
	var kept []string
	for _, p := range parts {
		if p != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, "\n\n")
}
