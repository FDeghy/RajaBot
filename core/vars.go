package core

import (
	"errors"
	"sync"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	workers       map[Work]chan struct{}
	mutex         *sync.RWMutex
	res           chan *TrainData
	Bot           *gotgbot.Bot
	stations      *raja.Stations
	userTimeCache map[userCache]int64
)

var (
	ErrTrainNotFound    = errors.New("درخواست یافت نشد")
	ErrTrainAlreadyDone = errors.New("درخواست قبلا لغو شده")
)
