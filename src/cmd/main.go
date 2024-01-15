package main

import (
	"comm-inv-poc/src/internal/configs"
	"comm-inv-poc/src/internal/core/manager"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	cfg, err := configs.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting %s", cfg.App.Name)

	manager.Start()
}
