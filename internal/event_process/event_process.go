package eventprocess

import (
	"fmt"
	
	"biathlon_events_parser/internal/models"
)

func ProcessEvents(events []*models.Event, cfg *models.Config) map[int]*models.Competitor {
	competitors := make(map[int]*models.Competitor)

	for _, ev := range events {
		id := ev.CompetitorID
		if _, exists := competitors[id]; !exists {
			competitors[id] = &models.Competitor{
				ID:       id,
				ShotsHit: 5,
			}
		}
		comp := competitors[id]

		switch ev.EventID {
		case 1:
			fmt.Printf("[%s] Competitor(%d) has registered\n",
				ev.Time.Format("15:04:05.000"), id)
		case 2:
			comp.ScheduledStart = ev.StartTime
			fmt.Printf("[%s] Scheduled start time for Competitor(%d) is %s (by draw)\n",
				ev.Time.Format("15:04:05.000"), id, ev.StartTime.Format("15:04:05.000"))
		case 3:
			fmt.Printf("[%s] Competitor(%d) is on the start line\n",
				ev.Time.Format("15:04:05.000"), id)
		case 4:
			comp.ActualStart = ev.Time
			comp.LastLapEnd = ev.Time
			fmt.Printf("[%s] Competitor(%d) has started\n",
				ev.Time.Format("15:04:05.000"), id)
		case 5:
			fmt.Printf("[%s] Competitor(%d) entered the firing range(%d)\n",
				ev.Time.Format("15:04:05.000"), id, ev.FiringRange)
		case 6:
			comp.Hits++
			fmt.Printf("[%s] Target(%d) was hit by Competitor(%d)\n",
				ev.Time.Format("15:04:05.000"), ev.Target, id)
		case 7:
			fmt.Printf("[%s] Competitor(%d) left the firing range\n",
				ev.Time.Format("15:04:05.000"), id)
		case 8:
			fmt.Printf("[%s] Competitor(%d) entered the penalty lap\n",
				ev.Time.Format("15:04:05.000"), id)
			comp.LastLapEnd = ev.Time
		case 9:
			if !comp.LastLapEnd.IsZero() {
				penalty := ev.Time.Sub(comp.LastLapEnd)
				comp.PenaltyTime += penalty
				if comp.PenaltyTime > 0 {
					comp.PenaltySpeed = float64(cfg.PenaltyLen) / comp.PenaltyTime.Seconds()
				}

				fmt.Printf("[%s] Competitor(%d) left the penalty lap\n",
					ev.Time.Format("15:04:05.000"), id)
			}
		case 10:
			if !comp.LastLapEnd.IsZero() {
				lapTime := ev.Time.Sub(comp.LastLapEnd)
				comp.LapTimes = append(comp.LapTimes, lapTime)
				speed := float64(cfg.LapLen) / lapTime.Seconds()
				comp.LapSpeeds = append(comp.LapSpeeds, speed)
				comp.LastLapEnd = ev.Time
				comp.Status = "[Finished]"

				fmt.Printf("[%s] Competitor(%d) finished a lap\n",
					ev.Time.Format("15:04:05.000"), id)
			}
		case 11:
			comp.Status = "[NotFinished]"
			comp.Comment = ev.Comment
			fmt.Printf("[%s] Competitor(%d) cannot continue: %s\n",
				ev.Time.Format("15:04:05.000"), id, ev.Comment)
		}
	}
	return competitors
}
