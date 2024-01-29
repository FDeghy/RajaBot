package bot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var (
	bot        *gotgbot.Bot
	updater    *ext.Updater
	dispatcher *ext.Dispatcher
)
