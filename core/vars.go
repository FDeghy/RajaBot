package core

import (
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	workers map[Work]chan struct{}
	mutex   sync.RWMutex
	res     chan *TrainData
	Bot     *gotgbot.Bot
)
