package checks_test

import (
	"strings"
	"testing"

	"github.com/djadmin/fort/internal/checks"
)

// ── Registry ────────────────────────────────────────────────────────────────

func TestRegistryCount(t *testing.T) {
	got := len(checks.All())
	const want = 16
	if got != want {
		t.Errorf("All() returned %d checks, want %d — update this test when adding a check", got, want)
	}
}

func TestNoDuplicateIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range checks.All() {
		if seen[c.ID()] {
			t.Errorf("duplicate check ID: %q", c.ID())
		}
		seen[c.ID()] = true
	}
}

// ── Interface contract ───────────────────────────────────────────────────────

func TestAllChecksHaveMetadata(t *testing.T) {
	for _, c := range checks.All() {
		if c.ID() == "" {
			t.Errorf("%T: empty ID", c)
		}
		if c.Name() == "" {
			t.Errorf("%s: empty Name", c.ID())
		}
	}
}

func TestAllChecksRun(t *testing.T) {
	validStatuses := map[checks.Status]bool{
		checks.StatusPass: true,
		checks.StatusFail: true,
		checks.StatusWarn: true,
	}
	for _, c := range checks.All() {
		r := c.Run()
		if r.ID != c.ID() {
			t.Errorf("%s: result ID %q != check ID", c.ID(), r.ID)
		}
		if r.Name != c.Name() {
			t.Errorf("%s: result Name %q != check Name", c.ID(), r.Name)
		}
		if !validStatuses[r.Status] {
			t.Errorf("%s: unknown status %q", c.ID(), r.Status)
		}
		if r.Current == "" {
			t.Errorf("%s: empty Current value", c.ID())
		}
		if r.Expected == "" {
			t.Errorf("%s: empty Expected value", c.ID())
		}
	}
}

// TestFixableConsistency is the contract that keeps the CI honest:
// Result.Fixable must always equal Check.Fixable() regardless of which
// code path Run() takes. If this breaks, a branch in Run() forgot to set Fixable.
func TestFixableConsistency(t *testing.T) {
	for _, c := range checks.All() {
		r := c.Run()
		if r.Fixable != c.Fixable() {
			t.Errorf("%s: Result.Fixable=%v but Check.Fixable()=%v — all Run() branches must set Fixable: c.Fixable()",
				c.ID(), r.Fixable, c.Fixable())
		}
	}
}

func TestFixableChecksHaveDescription(t *testing.T) {
	for _, c := range checks.All() {
		if c.Fixable() && c.FixDescription() == "" {
			t.Errorf("%s: Fixable() is true but FixDescription() is empty", c.ID())
		}
		if !c.Fixable() && c.FixDescription() != "" {
			t.Errorf("%s: Fixable() is false but FixDescription() is non-empty", c.ID())
		}
	}
}

// ── Framework mappings ───────────────────────────────────────────────────────

func TestAllChecksHaveFrameworkMappings(t *testing.T) {
	for _, c := range checks.All() {
		fw := checks.FrameworksFor(c.ID())
		if len(fw) == 0 {
			t.Errorf("%s: no framework mappings — add to frameworks.go", c.ID())
		}
	}
}

func TestAllChecksHaveSOC2Mapping(t *testing.T) {
	for _, c := range checks.All() {
		hasSOC2 := false
		for _, f := range checks.FrameworksFor(c.ID()) {
			if f.Name == "SOC 2" {
				hasSOC2 = true
				break
			}
		}
		if !hasSOC2 {
			t.Errorf("%s: missing SOC 2 mapping", c.ID())
		}
	}
}

func TestFrameworkControlsNonEmpty(t *testing.T) {
	for _, c := range checks.All() {
		for _, fw := range checks.FrameworksFor(c.ID()) {
			if fw.Name == "" {
				t.Errorf("%s: framework entry has empty Name", c.ID())
			}
			if len(fw.Controls) == 0 {
				t.Errorf("%s / %s: no controls listed", c.ID(), fw.Name)
			}
			for _, ctrl := range fw.Controls {
				if ctrl == "" {
					t.Errorf("%s / %s: empty control string", c.ID(), fw.Name)
				}
			}
		}
	}
}

func TestFrameworksForUnknownID(t *testing.T) {
	fw := checks.FrameworksFor("nonexistent-check-id")
	if fw != nil {
		t.Errorf("FrameworksFor(unknown) = %v, want nil", fw)
	}
}

// TestKnownFrameworks verifies the four expected frameworks are present.
func TestKnownFrameworks(t *testing.T) {
	want := []string{"SOC 2", "ISO 27001", "NIST CSF", "CIS v8"}
	for _, c := range checks.All() {
		fwNames := map[string]bool{}
		for _, f := range checks.FrameworksFor(c.ID()) {
			fwNames[f.Name] = true
		}
		for _, w := range want {
			if !fwNames[w] {
				t.Errorf("%s: missing framework %q", c.ID(), w)
			}
		}
	}
}

// ── Status semantics ─────────────────────────────────────────────────────────

// TestStatusConstants guards against typos in status string values.
func TestStatusConstants(t *testing.T) {
	if string(checks.StatusPass) != "pass" {
		t.Errorf("StatusPass = %q, want %q", checks.StatusPass, "pass")
	}
	if string(checks.StatusFail) != "fail" {
		t.Errorf("StatusFail = %q, want %q", checks.StatusFail, "fail")
	}
	if string(checks.StatusWarn) != "warn" {
		t.Errorf("StatusWarn = %q, want %q", checks.StatusWarn, "warn")
	}
}

// TestCheckIDsAreLowercase ensures IDs are safe as JSON keys and CLI args.
func TestCheckIDsAreLowercase(t *testing.T) {
	for _, c := range checks.All() {
		if c.ID() != strings.ToLower(c.ID()) {
			t.Errorf("%s: ID must be lowercase", c.ID())
		}
		if strings.ContainsAny(c.ID(), " \t\n") {
			t.Errorf("%s: ID must not contain whitespace", c.ID())
		}
	}
}
