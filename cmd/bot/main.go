package main

import (
	"go_ctl_bot/internal/bot"
	"go_ctl_bot/internal/config"
	"log"
)

func main() {
	cfgPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	bot.RunBot(cfg)
}
