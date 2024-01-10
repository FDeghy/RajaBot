package core

import (
	"RajaBot/config"
	"errors"
	"sync"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	Bot      *gotgbot.Bot
	stations *raja.Stations
)

var (
	mutex         = &sync.RWMutex{}
	workers       = make(map[Work]chan struct{}) //fetchWorkers
	userTimeCache = make(map[userCache]int64)
	res           = make(chan *TrainData, config.Cfg.Raja.Buffer)
)

var (
	ErrTrainNotFound    = errors.New("درخواست یافت نشد")
	ErrTrainAlreadyDone = errors.New("درخواست قبلا لغو شده")
)
