package report

import (
	"biathlon_events_parser/internal/models"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLapTimesTotal(t *testing.T) {
	comp := &models.Competitor{
		LapTimes: []time.Duration{
			30 * time.Second,
			45 * time.Second,
			60 * time.Second,
		},
	}

	total := lapTimesTotal(comp)
	expected := 135 * time.Second

	if total != expected {
		t.Errorf("lapTimesTotal() = %v, want %v", total, expected)
	}

	emptyComp := &models.Competitor{
		LapTimes: []time.Duration{},
	}

	total = lapTimesTotal(emptyComp)
	if total != 0 {
		t.Errorf("lapTimesTotal() with empty array = %v, want 0", total)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "Minutes and seconds",
			duration: 2*time.Minute + 35*time.Second + 500*time.Millisecond,
			want:     "02:35.500",
		},
		{
			name:     "Only seconds",
			duration: 45*time.Second + 100*time.Millisecond,
			want:     "00:45.100",
		},
		{
			name:     "With hours",
			duration: 1*time.Hour + 23*time.Minute + 45*time.Second + 678*time.Millisecond,
			want:     "01:23:45.678",
		},
		{
			name:     "Negative duration",
			duration: -30 * time.Second,
			want:     "00:30.000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %v, want %v", tt.duration, got, tt.want)
			}
		})
	}
}

func TestUpdateCompetitorStatus(t *testing.T) {
	tests := []struct {
		name       string
		competitor *models.Competitor
		wantStatus string
	}{
		{
			name: "Already NotFinished status",
			competitor: &models.Competitor{
				Status: "[NotFinished]",
			},
			wantStatus: "[NotFinished]",
		},
		{
			name: "Has started but not finished",
			competitor: &models.Competitor{
				ActualStart: time.Now(),
				Status:      "",
			},
			wantStatus: "[Finished]",
		},
		{
			name: "Not started but scheduled",
			competitor: &models.Competitor{
				ScheduledStart: time.Now(),
				ActualStart:    time.Time{}, 
				Status:         "",
			},
			wantStatus: "[NotStarted]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCompetitorStatus(tt.competitor)
			if tt.competitor.Status != tt.wantStatus {
				t.Errorf("updateCompetitorStatus() got status = %v, want %v",
					tt.competitor.Status, tt.wantStatus)
			}
		})
	}
}

func TestCalculateTotalTime(t *testing.T) {
	comp1 := &models.Competitor{
		LapTimes: []time.Duration{
			1 * time.Minute,
			2 * time.Minute,
		},
	}

	expected1 := 3 * time.Minute
	if got := calculateTotalTime(comp1); got != expected1 {
		t.Errorf("calculateTotalTime() = %v, want %v", got, expected1)
	}

	comp2 := &models.Competitor{
		LapTimes:    []time.Duration{1 * time.Minute},
		PenaltyTime: 30 * time.Second,
	}

	expected2 := 1*time.Minute + 30*time.Second
	if got := calculateTotalTime(comp2); got != expected2 {
		t.Errorf("calculateTotalTime() = %v, want %v", got, expected2)
	}

	now := time.Now()
	comp3 := &models.Competitor{
		LapTimes:       []time.Duration{1 * time.Minute},
		ScheduledStart: now,
		ActualStart:    now.Add(20 * time.Second),
	}

	expected3 := 1*time.Minute + 20*time.Second
	if got := calculateTotalTime(comp3); got != expected3 {
		t.Errorf("calculateTotalTime() = %v, want %v", got, expected3)
	}
}

func TestMakeReport(t *testing.T) {
	cfg := &models.Config{
		Laps:       2,
		LapLen:     1000,
		PenaltyLen: 150,
	}

	competitors := map[int]*models.Competitor{
		1: {
			ID:          1,
			ActualStart: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
			LapTimes:    []time.Duration{1 * time.Minute, 1*time.Minute + 10*time.Second},
			LapSpeeds:   []float64{16.67, 15.38},
			Status:      "[Finished]",
			Hits:        4,
			ShotsHit:    5,
		},
		2: {
			ID:           2,
			ActualStart:  time.Date(2023, 1, 1, 10, 1, 0, 0, time.UTC),
			LapTimes:     []time.Duration{55 * time.Second, 1 * time.Minute},
			LapSpeeds:    []float64{18.18, 16.67},
			PenaltyTime:  15 * time.Second,
			PenaltySpeed: 10.0,
			Status:       "[Finished]",
			Hits:         3,
			ShotsHit:     5,
		},
		3: {
			ID:          3,
			ActualStart: time.Date(2023, 1, 1, 10, 2, 0, 0, time.UTC),
			LapTimes:    []time.Duration{50 * time.Second},
			LapSpeeds:   []float64{20.0},
			Status:      "[NotFinished]",
			Comment:     "Fell",
			Hits:        2,
			ShotsHit:    5,
		},
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := MakeReport(competitors, cfg)
	if err != nil {
		t.Errorf("MakeReport() returned error: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "[Finished]") ||
		!strings.Contains(output, "[NotFinished]") {
		t.Errorf("MakeReport() output doesn't contain required status markers")
	}

	for _, id := range []int{1, 2, 3} {
		if !strings.Contains(output, fmt.Sprintf("%d", id)) {
			t.Errorf("MakeReport() output doesn't contain competitor ID %d", id)
		}
	}
}
