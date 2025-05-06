package eventprocess

import (
	"testing"
	"time"

	"biathlon_events_parser/internal/models"
)

func newTestEvent(timeStr string, eventID, compID int) *models.Event {
	t, _ := time.Parse("15:04:05.000", timeStr)
	return &models.Event{
		Time:         t,
		EventID:      eventID,
		CompetitorID: compID,
	}
}

func TestProcessEvents(t *testing.T) {
	cfg := &models.Config{
		Laps:       3,
		LapLen:     1000,
		PenaltyLen: 150,
	}

	events := []*models.Event{
		newTestEvent("10:00:00.000", 1, 10),
		newTestEvent("10:01:00.000", 2, 10),
		newTestEvent("10:02:00.000", 3, 10),
		newTestEvent("10:03:00.000", 4, 10),
		newTestEvent("10:04:00.000", 5, 10),
		newTestEvent("10:04:10.000", 6, 10),
		newTestEvent("10:04:20.000", 7, 10),
		newTestEvent("10:05:00.000", 10, 10),
	}

	competitors := ProcessEvents(events, cfg)

	if len(competitors) != 1 {
		t.Errorf("ProcessEvents() returned %d competitors, want 1", len(competitors))
	}

	comp := competitors[10]
	if comp == nil {
		t.Fatal("Competitor 10 not found in results")
	}

	if comp.ID != 10 {
		t.Errorf("Competitor ID = %d, want 10", comp.ID)
	}

	expectedStart, _ := time.Parse("15:04:05.000", "10:03:00.000")
	if !comp.ActualStart.Equal(expectedStart) {
		t.Errorf("ActualStart = %v, want %v", comp.ActualStart, expectedStart)
	}

	if len(comp.LapTimes) != 1 {
		t.Errorf("LapTimes count = %d, want 1", len(comp.LapTimes))
	}

	expectedLapTime := 2 * time.Minute
	if comp.LapTimes[0] != expectedLapTime {
		t.Errorf("LapTime = %v, want %v", comp.LapTimes[0], expectedLapTime)
	}

	expectedSpeed := float64(cfg.LapLen) / expectedLapTime.Seconds()
	if comp.LapSpeeds[0] != expectedSpeed {
		t.Errorf("LapSpeed = %v, want %v", comp.LapSpeeds[0], expectedSpeed)
	}

	if comp.Hits != 1 {
		t.Errorf("Hits = %d, want 1", comp.Hits)
	}

	if comp.Status != "[Finished]" {
		t.Errorf("Status = %s, want [Finished]", comp.Status)
	}
}

func TestProcessEventsMultipleCompetitors(t *testing.T) {
	cfg := &models.Config{
		Laps:       3,
		LapLen:     1000,
		PenaltyLen: 150,
	}

	events := []*models.Event{
		newTestEvent("10:00:00.000", 1, 10),
		newTestEvent("10:01:00.000", 4, 10),

		newTestEvent("10:00:30.000", 1, 20),
		newTestEvent("10:01:30.000", 4, 20),
		newTestEvent("10:02:00.000", 10, 10),
		newTestEvent("10:02:30.000", 10, 20),
	}

	competitors := ProcessEvents(events, cfg)

	if len(competitors) != 2 {
		t.Errorf("ProcessEvents() returned %d competitors, want 2", len(competitors))
	}

	comp10 := competitors[10]
	comp20 := competitors[20]

	if comp10 == nil || comp20 == nil {
		t.Fatal("One of the competitors not found in results")
	}

	expected10 := 1 * time.Minute
	if comp10.LapTimes[0] != expected10 {
		t.Errorf("Competitor 10 LapTime = %v, want %v", comp10.LapTimes[0], expected10)
	}

	expected20 := 1 * time.Minute
	if comp20.LapTimes[0] != expected20 {
		t.Errorf("Competitor 20 LapTime = %v, want %v", comp20.LapTimes[0], expected20)
	}
}

func TestProcessEventsPenalty(t *testing.T) {
	cfg := &models.Config{
		Laps:       3,
		LapLen:     1000,
		PenaltyLen: 150,
	}

	events := []*models.Event{
		newTestEvent("10:00:00.000", 1, 10),
		newTestEvent("10:01:00.000", 4, 10),
		newTestEvent("10:02:00.000", 8, 10),
		newTestEvent("10:02:30.000", 9, 10),
		newTestEvent("10:03:00.000", 10, 10),
	}

	competitors := ProcessEvents(events, cfg)
	comp := competitors[10]

	expectedPenalty := 30 * time.Second
	if comp.PenaltyTime != expectedPenalty {
		t.Errorf("PenaltyTime = %v, want %v", comp.PenaltyTime, expectedPenalty)
	}

	expectedPenaltySpeed := float64(cfg.PenaltyLen) / expectedPenalty.Seconds()
	if comp.PenaltySpeed != expectedPenaltySpeed {
		t.Errorf("PenaltySpeed = %v, want %v", comp.PenaltySpeed, expectedPenaltySpeed)
	}
}
