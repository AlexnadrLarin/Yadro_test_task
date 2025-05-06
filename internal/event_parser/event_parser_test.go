package eventparser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEventLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantID   int
		wantComp int
		wantTime string
	}{
		{
			name:     "Success registration test  регистрации",
			input:    "[09:30:00.000] 1 10",
			wantErr:  false,
			wantID:   1,
			wantComp: 10,
			wantTime: "09:30:00.000",
		},
		{
			name:     "Success scheduled start test",
			input:    "[09:30:00.000] 2 10 10:00:00.000",
			wantErr:  false,
			wantID:   2,
			wantComp: 10,
			wantTime: "09:30:00.000",
		},
		{
			name:    "Invalid format test",
			input:   "09:30:00.000 1 10",
			wantErr: true,
		},
		{
			name:    "Invalid time format test",
			input:   "[09:30:00] 1 10",
			wantErr: true,
		},
		{
			name:    "Invalid Event ID test",
			input:   "[09:30:00.000] X 10",
			wantErr: true,
		},
		{
			name:    "Invalid competitor ID test",
			input:   "[09:30:00.000] 1 X",
			wantErr: true,
		},
		{
			name:    "Not enough fields test",
			input:   "[09:30:00.000] 1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseEventLine(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseEventLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if got.EventID != tt.wantID {
				t.Errorf("parseEventLine() got EventID = %v, want %v", got.EventID, tt.wantID)
			}
			if got.CompetitorID != tt.wantComp {
				t.Errorf("parseEventLine() got CompetitorID = %v, want %v", got.CompetitorID, tt.wantComp)
			}

			gotTime := got.Time.Format("15:04:05.000")
			if gotTime != tt.wantTime {
				t.Errorf("parseEventLine() got Time = %v, want %v", gotTime, tt.wantTime)
			}
		})
	}
}

func TestParseEvents(t *testing.T) {
	content := `[09:30:00.000] 1 10
[09:30:05.000] 2 10 10:00:00.000
[09:30:10.000] 3 10
[09:30:15.000] 4 10`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_events.txt")

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current dir: %v", err)
	}
	defer os.Chdir(origWd)

	os.Chdir(tmpDir)

	events, err := ParseEvents("/test_events.txt")
	if err != nil {
		t.Fatalf("ParseEvents() error = %v", err)
	}

	if len(events) != 4 {
		t.Errorf("ParseEvents() returned %d events, want 4", len(events))
	}

	wantTimes := []string{
		"09:30:00.000",
		"09:30:05.000",
		"09:30:10.000",
		"09:30:15.000",
	}

	for i, event := range events {
		if event.Time.Format("15:04:05.000") != wantTimes[i] {
			t.Errorf("Event %d has wrong time: got %s, want %s",
				i, event.Time.Format("15:04:05.000"), wantTimes[i])
		}
	}
}

func TestParseEventsFileErrors(t *testing.T) {
	_, err := ParseEvents("/does_not_exist.txt")
	if err == nil {
		t.Error("ParseEvents() should return error for non-existent file")
	}
}
