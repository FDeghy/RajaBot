package main

import (
	"RajaBot/bot"
	"RajaBot/bot/adminHandler"
	"RajaBot/config"
	"RajaBot/core"
	"RajaBot/database"
	"RajaBot/payment"
	"RajaBot/prometheus"
	"RajaBot/tools"
	"flag"
	"fmt"
	"log"
)

func main() {
	fl := flag.Bool("addall", false, "add free sub to all users that have sub from 1 farvardin")
	flag.Parse()

	if *fl {
		addall()
		return
	}

	err := config.ParseConfig("config.toml")
	if err != nil {
		log.Fatalln("failed to parse config.")
	}

	err = bot.CreateBot()
	if err != nil {
		log.Fatalf("failed to create bot.\nerror: %v\n", err)
	}

	err = database.StartDatabase()
	if err != nil {
		log.Fatalln("failed to start database.")
	}
	defer database.CloseDatabase()

	err = core.StartCore()
	if err != nil {
		log.Fatalln("failed to start core.")
	}

	err = prometheus.StartProm()
	if err != nil {
		log.Fatalf("failed to start prometheus.\nerror: %v\n", err)
	}

	err = payment.StartPaymentServer()
	if err != nil {
		log.Fatalf("failed to start payment server.\nerror: %v\n", err)
	}

	err = bot.StartBot()
	if err != nil {
		log.Fatalf("failed to start bot.\nerror: %v\n", err)
	}
}

func addall() {
	days := 30

	err := config.ParseConfig("config.toml")
	if err != nil {
		log.Fatalln("failed to parse config.")
	}

	err = bot.CreateBot()
	if err != nil {
		log.Fatalf("failed to create bot.\nerror: %v\n", err)
	}

	err = database.StartDatabase()
	if err != nil {
		log.Fatalln("failed to start database.")
	}

	for _, user := range database.GetAllSubscription() {
		if (user.ExpirationDate-2592000) >= 1710361800 && database.IsHavePayment(user.UserID) {
			tools.AddDaysSub(user.UserID, days)
			_, err := tools.Bot.SendMessage(user.UserID, fmt.Sprintf(adminHandler.AddSubMsg, days), nil)
			if err != nil {
				log.Printf("Error: %v, User: %v", err, user.UserID)
				continue
			}
			log.Printf("add %v days to %v", days, user.UserID)
		}
	}
}
