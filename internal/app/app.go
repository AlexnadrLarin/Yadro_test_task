package app

import (
	"biathlon_events_parser/internal/config"
	"log"
)

func Run()  {
	cfg, err := config.LoadConfig("/configs/config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)
}