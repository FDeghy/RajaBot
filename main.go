package main

import (
	"RajaBot/bot"
	"RajaBot/config"
	"RajaBot/core"
	"RajaBot/database"
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

	err = bot.StartBot()
	if err != nil {
		log.Fatalln(err)
	}
}
