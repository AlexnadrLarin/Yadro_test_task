package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(origWd)

	os.Chdir(tmpDir)

	configsDir := filepath.Join(tmpDir, "configs")
	if err := os.Mkdir(configsDir, 0755); err != nil {
		t.Fatalf("Failed to create configs dir: %v", err)
	}

	configContent := `{
		"laps": 3,
		"lapLen": 1000,
		"penaltyLen": 150,
		"firingLines": 5,
		"start": "10:00:00.000",
		"startDelta": "00:30:00"
	}`

	configFile := filepath.Join(configsDir, "config.json")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadConfig("/configs/config.json")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Laps != 3 {
		t.Errorf("cfg.Laps = %d, want 3", cfg.Laps)
	}

	if cfg.LapLen != 1000 {
		t.Errorf("cfg.LapLen = %d, want 1000", cfg.LapLen)
	}

	if cfg.PenaltyLen != 150 {
		t.Errorf("cfg.PenaltyLen = %d, want 150", cfg.PenaltyLen)
	}

	if cfg.FiringLines != 5 {
		t.Errorf("cfg.FiringLines = %d, want 5", cfg.FiringLines)
	}

	expectedStartTime, _ := time.Parse("15:04:05.000", "10:00:00.000")
	if !cfg.StartTime.Equal(expectedStartTime) {
		t.Errorf("cfg.StartTime = %v, want %v", cfg.StartTime, expectedStartTime)
	}

	expectedDelta := 30 * time.Minute
	if cfg.StartDelta != expectedDelta {
		t.Errorf("cfg.StartDelta = %v, want %v", cfg.StartDelta, expectedDelta)
	}
}

func TestLoadConfigErrors(t *testing.T) {
	_, err := LoadConfig("/non_existent_config.json")
	if err == nil {
		t.Error("LoadConfig() should return error for non-existent file")
	}

	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	configsDir := filepath.Join(tmpDir, "configs")
	if err := os.Mkdir(configsDir, 0755); err != nil {
		t.Fatalf("Failed to create configs dir: %v", err)
	}

	invalidContent := `{ "laps": 3, invalid json }`
	invalidFile := filepath.Join(configsDir, "invalid.json")
	if err := os.WriteFile(invalidFile, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	_, err = LoadConfig("/configs/invalid.json")
	if err == nil {
		t.Error("LoadConfig() should return error for invalid JSON")
	}

	invalidTimeContent := `{
		"laps": 3,
		"lapLen": 1000,
		"penaltyLen": 150,
		"firingLines": 5,
		"start": "invalid time",
		"startDelta": "00:30:00"
	}`

	invalidTimeFile := filepath.Join(configsDir, "invalid_time.json")
	if err := os.WriteFile(invalidTimeFile, []byte(invalidTimeContent), 0644); err != nil {
		t.Fatalf("Failed to write invalid time config file: %v", err)
	}

	_, err = LoadConfig("/configs/invalid_time.json")
	if err == nil {
		t.Error("LoadConfig() should return error for invalid time format")
	}
}
