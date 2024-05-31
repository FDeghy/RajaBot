package core

import (
	"RajaBot/config"
	"RajaBot/database"
	siteapi "RajaBot/siteApi"
	"RajaBot/tools"
	"log"
	"slices"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
	ptime "github.com/yaa110/go-persian-calendar"
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

	go rtGetTrainsWorker()
	log.Printf("Core -> rtGetTrainsWorker started.")

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

func UpdateRtsTrains(route *siteapi.Route, date ptime.Time) ([]siteapi.Train, error) {
	trainList, err := siteapi.GetTrains(route.Src, route.Dst, date.Format("yyyy/MM/dd"))
	if err != nil || len(trainList) == 0 {
		return nil, err
	}

	oldTrainList := database.GetRtsByDate(route.Src, route.Dst, date)
	newTrainList := unionRtsData(oldTrainList, tools.SlicePtrToSlice(trainList))

	database.SetRtsByDate(route.Src, route.Dst, date, newTrainList)

	return newTrainList, nil
}
