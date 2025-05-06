package report

import (
	"biathlon_events_parser/internal/models"
	"fmt"
	"sort"
	"strings"
	"time"
)

func MakeReport(competitors map[int]*models.Competitor, cfg *models.Config) error {
	var list []*models.Competitor
	for _, comp := range competitors {
		updateCompetitorStatus(comp)
		list = append(list, comp)
	}

	sort.Slice(list, func(i, j int) bool {
		ci, cj := list[i], list[j]
		finishedI := ci.Status == "[Finished]" && !ci.ActualStart.IsZero()
		finishedJ := cj.Status == "[Finished]" && !cj.ActualStart.IsZero()
		notStartedI := ci.Status == "[NotStarted]"
		notStartedJ := cj.Status == "[NotStarted]"

		if finishedI && finishedJ {
			totalI := calculateTotalTime(ci)
			totalJ := calculateTotalTime(cj)
			return totalI < totalJ
		}

		if finishedI != finishedJ {
			return finishedI
		}

		if notStartedI != notStartedJ {
			return !notStartedI
		}

		return ci.ID < cj.ID
	})

	for _, comp := range list {
		lapsStrs := make([]string, 0, cfg.Laps)
		for i := 0; i < cfg.Laps; i++ {
			if i < len(comp.LapTimes) {
				t := comp.LapTimes[i]
				sp := comp.LapSpeeds[i]
				lapsStrs = append(lapsStrs, fmt.Sprintf("{%s, %.3f}", formatDuration(t), sp))
			} else {
				lapsStrs = append(lapsStrs, "{,}")
			}
		}
		lapsPart := strings.Join(lapsStrs, ", ")

		penaltyPart := "{,}"
		if comp.PenaltyTime > 0 {
			penaltyPart = fmt.Sprintf("{%s, %.3f}", formatDuration(comp.PenaltyTime), comp.PenaltySpeed)
		}

		targets := fmt.Sprintf("%d/%d", comp.Hits, comp.ShotsHit)

		fmt.Printf("%-12s %2d  [%s] %s  %s\n",
			comp.Status, comp.ID, lapsPart, penaltyPart, targets)
	}

	return nil
}

func updateCompetitorStatus(comp *models.Competitor) {
	if comp.Status == "[NotFinished]" {
		return
	}

	if !comp.ActualStart.IsZero() {
		if comp.Status == "" {
			comp.Status = "[Finished]"
		}
	} else if !comp.ScheduledStart.IsZero() {
		comp.Status = "[NotStarted]"
	}
}

func calculateTotalTime(c *models.Competitor) time.Duration {
	total := lapTimesTotal(c) + c.PenaltyTime

	if !c.ScheduledStart.IsZero() && !c.ActualStart.IsZero() {
		startDiff := c.ActualStart.Sub(c.ScheduledStart)
		if startDiff > 0 {
			total += startDiff
		}
	}

	return total
}

func lapTimesTotal(c *models.Competitor) time.Duration {
	var total time.Duration
	for _, t := range c.LapTimes {
		total += t
	}
	return total
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	ms := d - s*time.Second
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms/time.Millisecond)
	}
	return fmt.Sprintf("%02d:%02d.%03d", m, s, ms/time.Millisecond)
}
