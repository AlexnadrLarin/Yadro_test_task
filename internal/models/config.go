package models

import "time"

type Config struct {
	Laps        int           `json:"laps"`
	LapLen      int           `json:"lapLen"`
	PenaltyLen  int           `json:"penaltyLen"`
	FiringLines int           `json:"firingLines"`
	StartTime   time.Time     `json:"start"`
	StartDelta  time.Duration `json:"startDelta"`
}
