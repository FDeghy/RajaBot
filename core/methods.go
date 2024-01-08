package core

import (
	"RajaBot/database"
	"fmt"

	"github.com/FDeghy/RajaGo/raja"
	ptime "github.com/yaa110/go-persian-calendar"
)

func sendAlert(trainWR *database.TrainWR, train *raja.GoTrains) {
	stations, _ := raja.GetStations()
	src, _ := stations.GetPersianName(trainWR.Src)
	dst, _ := stations.GetPersianName(trainWR.Dst)
	Bot.SendMessage(
		trainWR.UserID,
		fmt.Sprintf(
			AlertMsg,
			train.ExitTime,
			ptime.Unix(trainWR.Day, 0).Format(TrainDate),
			src,
			dst,
			int(train.Counting),
		),
		nil,
	)
}

func expireWork(trainId int, exitTime string) {
	oldTrains := database.GetActiveTrainWRsByTrainId(trainId)
	for _, i := range *oldTrains {
		src, _ := stations.GetPersianName(i.Src)
		dst, _ := stations.GetPersianName(i.Dst)
		Bot.SendMessage(
			i.UserID,
			fmt.Sprintf(
				ExpireMsg,
				exitTime,
				ptime.Unix(i.Day, 0).Format(TrainDate),
				src,
				dst,
			),
			nil,
		)
		i.IsDone = true
		database.UpdateTrainWR(&i)
	}
}
