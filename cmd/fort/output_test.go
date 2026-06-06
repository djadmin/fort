package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/djadmin/fort/internal/checks"
)

// captureStdout runs f and returns whatever it wrote to stdout.
func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// ── tally ────────────────────────────────────────────────────────────────────

func TestTally(t *testing.T) {
	cases := []struct {
		name     string
		results  []checks.Result
		wantPass int
		wantFail int
		wantWarn int
	}{
		{"all pass", []checks.Result{{Status: checks.StatusPass}, {Status: checks.StatusPass}}, 2, 0, 0},
		{"mixed", []checks.Result{{Status: checks.StatusPass}, {Status: checks.StatusFail}, {Status: checks.StatusWarn}}, 1, 1, 1},
		{"empty", nil, 0, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pass, fail, warn := tally(tc.results)
			if pass != tc.wantPass || fail != tc.wantFail || warn != tc.wantWarn {
				t.Errorf("tally() = %d/%d/%d, want %d/%d/%d", pass, fail, warn, tc.wantPass, tc.wantFail, tc.wantWarn)
			}
		})
	}
}

// ── anyFixable ───────────────────────────────────────────────────────────────

func TestAnyFixable(t *testing.T) {
	cases := []struct {
		name    string
		results []checks.Result
		want    bool
	}{
		{"fixable fail", []checks.Result{{Status: checks.StatusFail, Fixable: true}}, true},
		{"fixable pass — ignored", []checks.Result{{Status: checks.StatusPass, Fixable: true}}, false},
		{"non-fixable fail", []checks.Result{{Status: checks.StatusFail, Fixable: false}}, false},
		{"fixable warn", []checks.Result{{Status: checks.StatusWarn, Fixable: true}}, true},
		{"empty", nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := anyFixable(tc.results); got != tc.want {
				t.Errorf("anyFixable() = %v, want %v", got, tc.want)
			}
		})
	}
}

// ── statusIcon ───────────────────────────────────────────────────────────────

func TestStatusIcon(t *testing.T) {
	cases := []struct {
		status   checks.Status
		wantIcon string
	}{
		{checks.StatusPass, "✓"},
		{checks.StatusFail, "✗"},
		{checks.StatusWarn, "~"},
	}
	for _, tc := range cases {
		icon, _ := statusIcon(tc.status)
		if icon != tc.wantIcon {
			t.Errorf("statusIcon(%q) = %q, want %q", tc.status, icon, tc.wantIcon)
		}
	}
}

// ── toJSONPolicies ───────────────────────────────────────────────────────────

func TestToJSONPolicies(t *testing.T) {
	results := []checks.Result{
		{ID: "filevault", Name: "Disk encryption", Status: checks.StatusPass},
		{ID: "screenlock", Name: "Screen lock", Status: checks.StatusFail},
		{ID: "nonexistent", Name: "Unknown check", Status: checks.StatusWarn},
	}
	policies := toJSONPolicies(results)

	if len(policies) != len(results) {
		t.Fatalf("toJSONPolicies() len = %d, want %d", len(policies), len(results))
	}
	if len(policies[0].Frameworks) == 0 {
		t.Error("filevault: expected framework mappings")
	}
	if _, ok := policies[0].Frameworks["SOC 2"]; !ok {
		t.Error("filevault: missing SOC 2 mapping")
	}
	if policies[1].Status != checks.StatusFail {
		t.Errorf("policy[1].Status = %q, want fail", policies[1].Status)
	}
	if policies[2].Frameworks == nil {
		t.Error("unknown check: Frameworks must be empty map, not nil")
	}
}

// ── printJSON ────────────────────────────────────────────────────────────────

func TestPrintJSON(t *testing.T) {
	results := []checks.Result{
		{ID: "filevault", Name: "Disk encryption", Status: checks.StatusPass, Current: "on", Expected: "on"},
		{ID: "screenlock", Name: "Screen lock", Status: checks.StatusFail, Current: "off", Expected: "immediate"},
	}
	out := captureStdout(func() {
		printJSON(results, "test-host", "SN123", "15.0")
	})

	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("printJSON output is not valid JSON: %v\n%s", err, out)
	}
	if payload["tool"] != "fort" {
		t.Errorf("tool = %v, want fort", payload["tool"])
	}
	if payload["hostname"] != "test-host" {
		t.Errorf("hostname = %v, want test-host", payload["hostname"])
	}
	policies, ok := payload["policies"].([]any)
	if !ok || len(policies) != 2 {
		t.Errorf("policies count = %d, want 2", len(policies))
	}
}

// ── printHuman ───────────────────────────────────────────────────────────────

func TestPrintHuman(t *testing.T) {
	results := []checks.Result{
		{ID: "filevault", Name: "Disk encryption", Status: checks.StatusPass, Current: "on", Expected: "on"},
		{ID: "screenlock", Name: "Screen lock", Status: checks.StatusFail, Current: "off", Expected: "immediate", Fixable: true},
		{ID: "antivirus", Name: "Antivirus", Status: checks.StatusWarn, Current: "XProtect only", Expected: "third-party AV"},
	}
	out := captureStdout(func() {
		printHuman(results, "test-host", "15.0")
	})

	for _, want := range []string{"Disk encryption", "Screen lock", "Antivirus", "1/3", "test-host", "15.0", "fort --fix"} {
		if !strings.Contains(out, want) {
			t.Errorf("printHuman output missing %q", want)
		}
	}
}

// ── version ──────────────────────────────────────────────────────────────────

func TestVersion(t *testing.T) {
	if version == "" {
		t.Error("version must not be empty")
	}
	// Format check: must be semver-ish (e.g. "0.1.1")
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		t.Errorf("version %q is not in MAJOR.MINOR.PATCH format", version)
	}
}

// ── run (dry-run) ─────────────────────────────────────────────────────────────
// run() calls real OS checks. We only test the dry-run path here since it is
// read-only and never prompts. The interactive --fix prompt requires a TTY and
// cannot be unit tested; verify it manually with fort --fix.

func TestRunDryRun(t *testing.T) {
	out := captureStdout(func() {
		code, err := run(false, false, true, false, false, "")
		if err != nil {
			t.Errorf("run(dryRun) error: %v", err)
		}
		if code != 0 {
			t.Errorf("run(dryRun) exit code = %d, want 0", code)
		}
	})
	if !strings.Contains(out, "dry run") {
		t.Errorf("dry-run output missing 'dry run', got: %q", out)
	}
}

func TestRunJSONOutput(t *testing.T) {
	out := captureStdout(func() {
		run(true, false, false, false, false, "") //nolint
	})
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("--json output is not valid JSON: %v", err)
	}
	if payload["tool"] != "fort" {
		t.Errorf("--json tool = %v, want fort", payload["tool"])
	}
}

func TestRunOnly(t *testing.T) {
	// --only with a single known ID should produce JSON with exactly one policy.
	out := captureStdout(func() {
		run(true, false, false, false, false, "filevault") //nolint
	})
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("--json --only output is not valid JSON: %v", err)
	}
	policies, ok := payload["policies"].([]any)
	if !ok {
		t.Fatal("policies field missing or wrong type")
	}
	if len(policies) != 1 {
		t.Errorf("--only filevault: expected 1 policy, got %d", len(policies))
	}
	p := policies[0].(map[string]any)
	if p["id"] != "filevault" {
		t.Errorf("--only filevault: policy id = %v, want filevault", p["id"])
	}
}

func TestRunOnlyUnknownID(t *testing.T) {
	// An unknown ID should warn to stderr and return an error (no matching checks).
	var stderr bytes.Buffer
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	captureStdout(func() {
		run(true, false, false, false, false, "nonexistent_check_id") //nolint
	})

	w.Close()
	os.Stderr = oldStderr
	io.Copy(&stderr, r)

	if !strings.Contains(stderr.String(), "nonexistent_check_id") {
		t.Errorf("expected warning about unknown ID in stderr, got: %q", stderr.String())
	}
}
