package core

import (
	siteapi "RajaBot/siteApi"
	"RajaBot/siteApi/mrbilit"

	"github.com/FDeghy/RajaGo/raja"
)

type Work struct {
	Src     int
	Dst     int
	Day     int64
	ThrdApp int
}

type TrainData struct {
	TrainList *raja.TrainList
	Work      Work
}

type RtTrainData struct {
	TrainList []*siteapi.Train
	Work      Work
}

type ThrdAppTrainData struct {
	TrainList []*mrbilit.Trains
	Work      Work
}

type userCache struct {
	tgUserId int64
	trainId  int
}
