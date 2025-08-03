package version

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()

	if version == "" {
		t.Error("Version should not be empty")
	}

	// Version should be in a reasonable format
	if len(version) < 3 {
		t.Errorf("Version seems too short: %s", version)
	}

	t.Logf("Current version: %s", version)
}
