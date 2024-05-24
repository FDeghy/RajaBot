package core

import (
	siteapi "RajaBot/siteApi"

	"github.com/FDeghy/RajaGo/raja"
)

type Work struct {
	Src int
	Dst int
	Day int64
}

type TrainData struct {
	TrainList *raja.TrainList
	Work      Work
}

type RtTrainData struct {
	TrainList []*siteapi.Train
	Work      Work
}

type userCache struct {
	tgUserId int64
	trainId  int
}
