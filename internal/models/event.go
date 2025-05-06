package models

import "time"

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	StartTime    time.Time
	FiringRange  int
	Target       int
	Comment      string
}
