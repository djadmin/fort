package checks_test

import (
	"testing"

	"github.com/djadmin/fort/internal/checks"
)

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
	for _, c := range checks.All() {
		r := c.Run()
		if r.ID != c.ID() {
			t.Errorf("%s: result ID %q != check ID", c.ID(), r.ID)
		}
		if r.Status == "" {
			t.Errorf("%s: empty Status", c.ID())
		}
		if r.Current == "" {
			t.Errorf("%s: empty Current", c.ID())
		}
		if r.Expected == "" {
			t.Errorf("%s: empty Expected", c.ID())
		}
		validStatuses := map[checks.Status]bool{
			checks.StatusPass: true,
			checks.StatusFail: true,
			checks.StatusWarn: true,
		}
		if !validStatuses[r.Status] {
			t.Errorf("%s: unknown status %q", c.ID(), r.Status)
		}
	}
}

func TestFixableConsistency(t *testing.T) {
	for _, c := range checks.All() {
		r := c.Run()
		if r.Fixable != c.Fixable() {
			t.Errorf("%s: Result.Fixable=%v but Check.Fixable()=%v", c.ID(), r.Fixable, c.Fixable())
		}
		if c.Fixable() && c.FixDescription() == "" {
			t.Errorf("%s: Fixable() is true but FixDescription() is empty", c.ID())
		}
		if !c.Fixable() && c.FixDescription() != "" {
			t.Errorf("%s: Fixable() is false but FixDescription() is non-empty", c.ID())
		}
	}
}
