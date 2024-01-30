package core

import (
	"RajaBot/config"
	"RajaBot/database"
	"RajaBot/tools"
	"log"
	"slices"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

func StartCore() error {
	sts, err := raja.GetStations()
	if err != nil {
		return nil
	}
	stations = sts

	quit := make(chan struct{})
	for i := 0; i < config.Cfg.Raja.Worker; i++ {
		go procWorker(quit)
	}
	log.Printf("Core -> %v procWorker started.", config.Cfg.Raja.Worker)

	uncompTrainWRs := database.GetAllActiveTrainWRs()
	noHaveSubUsers := []int64{}
	// handle uncompleted tasks
	for _, i := range uncompTrainWRs {
		if !tools.CheckHaveSubscription(i.UserID) {
			if !slices.Contains(noHaveSubUsers, i.UserID) {
				noHaveSubUsers = append(noHaveSubUsers, i.UserID)
			}
			CancelWork(i.Id)
			continue
		}
		HandleGoFetch(i)
	}
	// handle users who have not sub
	for _, i := range noHaveSubUsers {
		sub := database.GetSubscription(i)
		if sub == nil {
			sub = database.NewSubscription(i)
			database.SaveSubscription(sub)
		}

		text, markup := tools.CreateSubStatus(sub)

		Bot.SendMessage(
			i,
			text,
			&gotgbot.SendMessageOpts{ReplyMarkup: markup},
		)
	}

	log.Printf("Core -> %v uncompleted task sent to handler.", len(uncompTrainWRs))

	return nil
}
