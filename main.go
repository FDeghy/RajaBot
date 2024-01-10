package main

import (
	"RajaBot/bot"
	"RajaBot/config"
	"RajaBot/core"
	"RajaBot/database"
	"RajaBot/prometheus"
	"log"
)

func main() {
	err := config.ParseConfig("config.toml")
	if err != nil {
		log.Fatalln("failed to parse config.")
	}

	err = database.StartDatabase()
	if err != nil {
		log.Fatalln("failed to start database.")
	}

	err = core.StartCore()
	if err != nil {
		log.Fatalln("failed to start core.")
	}

	err = prometheus.StartProm()
	if err != nil {
		log.Fatalf("failed to start prometheus.\nerror: %v\n", err)
	}

	err = bot.StartBot()
	if err != nil {
		log.Fatalf("failed to start bot.\nerror: %v\n", err)
	}
}
