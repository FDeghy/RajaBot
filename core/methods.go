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

func preCloseFetchWorker(wk Work) {

}
