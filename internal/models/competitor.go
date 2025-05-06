package models

import "time"

type Competitor struct {
	ID             int
	ScheduledStart time.Time   
	ActualStart    time.Time   
	LastLapEnd     time.Time   
	Hits           int         
	ShotsHit       int         
	PenaltyTime    time.Duration 
	PenaltySpeed   float64      
	LapTimes       []time.Duration 
	LapSpeeds      []float64       
	Status         string       
	Comment        string 
}