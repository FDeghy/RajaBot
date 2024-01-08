package core

import (
	"sync"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	workers  map[Work]chan struct{}
	mutex    sync.RWMutex
	res      chan *TrainData
	Bot      *gotgbot.Bot
	stations *raja.Stations
)
