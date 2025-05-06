package eventparser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"biathlon_events_parser/internal/models"
)

func ParseEvents(filePath string) ([]*models.Event, error) {
	homeDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var path strings.Builder
	path.WriteString(homeDir)
	path.WriteString(filePath)

	file, err := os.Open(path.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open event file: %w", err)
	}
	defer file.Close()

	var events []*models.Event
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		event, err := parseEventLine(line)
		if err != nil {
			fmt.Printf("Error parsing line %d: %v\n", lineNum, err)
			continue
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return events, nil
}

func parseEventLine(line string) (*models.Event, error) {
	parts := strings.SplitN(line, "]", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid event line format: %q", line)
	}

	timeStr := strings.Trim(parts[0], "[")
	t, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid event time format %q: %w", timeStr, err)
	}

	rest := strings.TrimSpace(parts[1])
	fields := strings.Fields(rest)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid number of event fields: %q", line)
	}

	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid event ID %q: %w", fields[0], err)
	}

	compID, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, fmt.Errorf("invalid competitor ID %q: %w", fields[1], err)
	}

	ev := models.Event{
		Time:         t,
		EventID:      id,
		CompetitorID: compID,
	}

	switch id {
	case 2:
		if len(fields) < 3 {
			return &ev, fmt.Errorf("missing start time parameter for event ID 2")
		}
		sched, err := time.Parse("15:04:05.000", fields[2])
		if err != nil {
			return &ev, fmt.Errorf("invalid start time format %q: %w", fields[2], err)
		}
		ev.StartTime = sched

	case 5:
		if len(fields) >= 3 {
			ev.FiringRange, err = strconv.Atoi(fields[2])
			if err != nil {
				return &ev, fmt.Errorf("invalid firing range %q: %w", fields[2], err)
			}
		}
	case 6:
		if len(fields) >= 3 {
			ev.Target, err = strconv.Atoi(fields[2])
			if err != nil {
				return &ev, fmt.Errorf("invalid target ID %q: %w", fields[2], err)
			}
		}
	case 11:
		if len(fields) >= 3 {
			ev.Comment = strings.Join(fields[2:], " ")
		}
	}

	return &ev, nil
}
