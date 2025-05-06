package app

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"biathlon_events_parser/internal/config"
	"biathlon_events_parser/internal/event_parser"
	"biathlon_events_parser/internal/event_process"
	"biathlon_events_parser/internal/report"
)

func Run() {
	var envPath strings.Builder

	homeDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	envPath.WriteString(homeDir)
	envPath.WriteString("/.env")

	err = godotenv.Load(envPath.String())
	if err != nil {
		log.Fatal(err)
	}

	configPath := os.Getenv("CONFIG_FILE_PATH")
	if configPath == "" {
		configPath = "/configs/config.json"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	eventsPath := os.Getenv("EVENTS_FILE_PATH")
	if eventsPath == "" {
		eventsPath = "/events/events"
	}

	events, err := eventparser.ParseEvents(eventsPath)
	if err != nil {
		log.Fatalf("failed to parse events: %v", err)
	}

	comp := eventprocess.ProcessEvents(events, cfg)

	err = report.MakeReport(comp, cfg)
	if err != nil {
		log.Fatalf("failed to create report: %v", err)
	}
}
