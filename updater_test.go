package updater

import (
	"os"
	"testing"
)

// Test that the install path is respecting the alread installed
// application
func TestExistingInstallPath(t *testing.T) {
	updater := NewUpdater("cmt", "1.0.0", "segersniels")
	path := updater.determineInstallPath()

	if path != "/usr/local/bin" {
		t.Fatalf("Expected /usr/local/bin, got %s", path)
	}
}

// Test that the install path is falling back to either GOBIN or /usr/local/bin
func TestFallbackInstallPath(t *testing.T) {
	updater := NewUpdater("foo", "1.0.0", "segersniels")

	os.Setenv("GOBIN", "/tmp")
	path := updater.determineInstallPath()
	if path != "/tmp" {
		t.Fatalf("Expected /tmp, got %s", path)
	}

	os.Unsetenv("GOBIN")
	path = updater.determineInstallPath()
	if path != "/usr/local/bin" {
		t.Fatalf("Expected /usr/local/bin, got %s", path)
	}
}
