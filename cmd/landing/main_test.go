package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanEmail(t *testing.T) {
	email, err := cleanEmail("  Person@Example.COM ")
	if err != nil {
		t.Fatalf("cleanEmail() error = %v", err)
	}
	if email != "person@example.com" {
		t.Fatalf("cleanEmail() = %q, want %q", email, "person@example.com")
	}
}

func TestCleanEmailRejectsInvalidAddress(t *testing.T) {
	if _, err := cleanEmail("not-an-email"); err == nil {
		t.Fatal("cleanEmail() accepted invalid address")
	}
}

func TestHandleWaitlistStoresSanitizedEmail(t *testing.T) {
	tmpDir := t.TempDir()
	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prevWD)

	req := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"email":" Person@Example.COM "}`))
	rec := httptest.NewRecorder()

	handleWaitlist(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, csvPath))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "person@example.com") {
		t.Fatalf("waitlist.csv = %q, want sanitized email", string(data))
	}
}

func TestHandleWaitlistRejectsInvalidEmail(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/waitlist", strings.NewReader(`{"email":"bad"}`))
	rec := httptest.NewRecorder()

	handleWaitlist(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
