package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"biathlon_events_parser/internal/models"
)

func LoadConfig(configPath string) (*models.Config, error) {
	homeDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var path strings.Builder
	path.WriteString(homeDir)
	path.WriteString(configPath)

	file, err := os.Open(path.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	temp := struct {
		Laps        int    `json:"laps"`
		LapLen      int    `json:"lapLen"`
		PenaltyLen  int    `json:"penaltyLen"`
		FiringLines int    `json:"firingLines"`
		Start       string `json:"start"`
		StartDelta  string `json:"startDelta"`
	}{}

	if err := dec.Decode(&temp); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	startTime, err := time.Parse("15:04:05.000", temp.Start)
	if err != nil {
		return nil, fmt.Errorf("invalid time format for Start in config: %w", err)
	}

	delta, err := time.Parse("15:04:05", temp.StartDelta)
	if err != nil {
		return nil, fmt.Errorf("invalid format for StartDelta: %s", temp.StartDelta)
	}

	startDelta := time.Duration(delta.Hour())*time.Hour +
		time.Duration(delta.Minute())*time.Minute +
		time.Duration(delta.Second())*time.Second

	cfg := models.Config{
		Laps:        temp.Laps,
		LapLen:      temp.LapLen,
		PenaltyLen:  temp.PenaltyLen,
		FiringLines: temp.FiringLines,
		StartTime:   startTime,
		StartDelta:  startDelta,
	}

	return &cfg, nil
}
