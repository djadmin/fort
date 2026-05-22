package main

import (
	"testing"

	"github.com/djadmin/fort/internal/checks"
)

func TestTally(t *testing.T) {
	cases := []struct {
		name        string
		results     []checks.Result
		wantPass    int
		wantFail    int
		wantWarn    int
	}{
		{
			name:     "all pass",
			results:  []checks.Result{{Status: checks.StatusPass}, {Status: checks.StatusPass}},
			wantPass: 2, wantFail: 0, wantWarn: 0,
		},
		{
			name:     "mixed",
			results:  []checks.Result{{Status: checks.StatusPass}, {Status: checks.StatusFail}, {Status: checks.StatusWarn}},
			wantPass: 1, wantFail: 1, wantWarn: 1,
		},
		{
			name:     "empty",
			results:  nil,
			wantPass: 0, wantFail: 0, wantWarn: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pass, fail, warn := tally(tc.results)
			if pass != tc.wantPass || fail != tc.wantFail || warn != tc.wantWarn {
				t.Errorf("tally() = %d/%d/%d, want %d/%d/%d",
					pass, fail, warn, tc.wantPass, tc.wantFail, tc.wantWarn)
			}
		})
	}
}

func TestAnyFixable(t *testing.T) {
	cases := []struct {
		name    string
		results []checks.Result
		want    bool
	}{
		{
			name:    "fixable fail",
			results: []checks.Result{{Status: checks.StatusFail, Fixable: true}},
			want:    true,
		},
		{
			name:    "fixable but already passing",
			results: []checks.Result{{Status: checks.StatusPass, Fixable: true}},
			want:    false,
		},
		{
			name:    "non-fixable fail",
			results: []checks.Result{{Status: checks.StatusFail, Fixable: false}},
			want:    false,
		},
		{
			name:    "fixable warn",
			results: []checks.Result{{Status: checks.StatusWarn, Fixable: true}},
			want:    true,
		},
		{
			name:    "empty",
			results: nil,
			want:    false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := anyFixable(tc.results)
			if got != tc.want {
				t.Errorf("anyFixable() = %v, want %v", got, tc.want)
			}
		})
	}
}

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

	// Known check must have framework mappings
	if len(policies[0].Frameworks) == 0 {
		t.Error("filevault policy: expected framework mappings, got none")
	}
	if _, ok := policies[0].Frameworks["SOC 2"]; !ok {
		t.Error("filevault policy: missing SOC 2 framework mapping")
	}
	if _, ok := policies[0].Frameworks["ISO 27001"]; !ok {
		t.Error("filevault policy: missing ISO 27001 framework mapping")
	}

	// Result fields must be preserved
	if policies[1].ID != "screenlock" {
		t.Errorf("policy[1].ID = %q, want %q", policies[1].ID, "screenlock")
	}
	if policies[1].Status != checks.StatusFail {
		t.Errorf("policy[1].Status = %q, want fail", policies[1].Status)
	}

	// Unknown check ID should produce empty (not nil-crash) frameworks
	if policies[2].Frameworks == nil {
		t.Error("unknown check policy: Frameworks should be empty map, not nil")
	}
}
