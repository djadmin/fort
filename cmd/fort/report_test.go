package main

import (
	"os"
	"strings"
	"testing"

	"github.com/djadmin/fort/internal/checks"
)

func TestWriteReport(t *testing.T) {
	results := []checks.Result{
		{
			ID: "passwordmgr", Name: "Password manager",
			Status: checks.StatusPass, Current: "1Password", Expected: "installed",
		},
		{
			ID: "filevault", Name: "Disk encryption",
			Status: checks.StatusFail, Current: "off", Expected: "on",
		},
		{
			ID: "antivirus", Name: "Antivirus / EDR",
			Status: checks.StatusWarn, Current: "XProtect only (built-in)", Expected: "third-party AV/EDR",
		},
	}

	tmp, err := os.CreateTemp("", "fort-report-test-*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.Close()

	if err := writeReport(results, "test-machine", "SNTEST123", "15.5", tmp.Name()); err != nil {
		t.Fatalf("writeReport() error: %v", err)
	}

	raw, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}
	html := string(raw)

	// Must be valid HTML
	if !strings.HasPrefix(strings.TrimSpace(html), "<!DOCTYPE html>") {
		t.Error("report does not start with <!DOCTYPE html>")
	}

	// Machine metadata
	mustContain(t, html, "test-machine", "machine hostname")
	mustContain(t, html, "SNTEST123", "serial number")
	mustContain(t, html, "15.5", "OS version")

	// Check names
	mustContain(t, html, "Password manager", "check name")
	mustContain(t, html, "Disk encryption", "check name")
	mustContain(t, html, "Antivirus / EDR", "check name")

	// Status values
	mustContain(t, html, "pass", "pass status")
	mustContain(t, html, "fail", "fail status")
	mustContain(t, html, "warn", "warn status")

	// Current values
	mustContain(t, html, "1Password", "current value")
	mustContain(t, html, "XProtect only (built-in)", "warn current value")

	// Expected shown for fail/warn, not for pass
	mustContain(t, html, "third-party AV/EDR", "warn expected value")

	// Framework mappings must appear
	mustContain(t, html, "SOC 2", "SOC 2 framework")
	mustContain(t, html, "ISO 27001", "ISO 27001 framework")
	mustContain(t, html, "NIST CSF", "NIST CSF framework")
	mustContain(t, html, "CIS v8", "CIS v8 framework")

	// Specific known control numbers
	mustContain(t, html, "CC6.1", "SOC 2 control")
	mustContain(t, html, "A.8.3", "ISO 27001 control")

	// Score
	mustContain(t, html, "1/3", "score")

	// fort branding
	mustContain(t, html, "fort", "brand")
	mustContain(t, html, "Confidential", "confidential label")
}

func TestWriteReportAllPass(t *testing.T) {
	results := []checks.Result{
		{ID: "filevault", Name: "Disk encryption", Status: checks.StatusPass, Current: "on", Expected: "on"},
	}
	tmp, err := os.CreateTemp("", "fort-report-allpass-*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.Close()

	if err := writeReport(results, "host", "SN1", "15.0", tmp.Name()); err != nil {
		t.Fatalf("writeReport() error: %v", err)
	}

	html := string(mustReadFile(t, tmp.Name()))
	mustContain(t, html, "1/1", "perfect score")
	mustContain(t, html, "s-pass", "green score class")
}

func TestWriteReportFixed(t *testing.T) {
	results := []checks.Result{
		{
			ID: "screenlock", Name: "Screen lock",
			Status: checks.StatusPass, Current: "immediate", Expected: "immediate",
			Fixed: true,
		},
	}
	tmp, err := os.CreateTemp("", "fort-report-fixed-*.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.Close()

	if err := writeReport(results, "host", "SN1", "15.0", tmp.Name()); err != nil {
		t.Fatalf("writeReport() error: %v", err)
	}
	html := string(mustReadFile(t, tmp.Name()))
	mustContain(t, html, "fixed", "fixed indicator")
}

func TestWriteReportInvalidPath(t *testing.T) {
	err := writeReport(nil, "h", "s", "1", "/nonexistent/dir/report.html")
	if err == nil {
		t.Error("writeReport() with invalid path should return error")
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func mustContain(t *testing.T, html, substr, label string) {
	t.Helper()
	if !strings.Contains(html, substr) {
		t.Errorf("report missing %s: %q not found in output", label, substr)
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
