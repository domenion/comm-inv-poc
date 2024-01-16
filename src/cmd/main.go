package main

import (
	"comm-inv-poc/src/internal/adapters"
	"comm-inv-poc/src/internal/configs"
	"comm-inv-poc/src/internal/core/app"
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
	tsm := adapters.NewTSMAdapter()
	cmt := adapters.NewCMTAdapter()
	a := app.New(tsm, cmt)
	a.ImportByItem()
}
