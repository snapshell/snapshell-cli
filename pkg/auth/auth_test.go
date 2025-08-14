package auth_test

import (
	"os"
	"testing"

	"github.com/snapshell/snapshell-cli/pkg/auth"
)

func TestGetConfigPath(t *testing.T) {
	path := auth.GetConfigPath()
	if path == "" {
		t.Error("Config path should not be empty")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Use a temp file for config
	tempFile := "test-snapshell-config.json"
	os.Setenv("HOME", ".") // force config to current dir
	defer os.Remove(tempFile)

	cfg := &auth.AuthConfig{Token: "abc123", APIUrl: "http://localhost:3000"}
	// Save config
	err := auth.SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
	// Load config
	loaded, err := auth.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loaded.Token != "abc123" || loaded.APIUrl != "http://localhost:3000" {
		t.Errorf("Loaded config does not match saved config: %+v", loaded)
	}
}
