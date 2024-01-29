package adminHandler

import (
	"RajaBot/tools"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func Load(d *ext.Dispatcher) {
	addsub := handlers.NewMessage(
		admin(message.HasPrefix("!add")),
		_addsub,
	)
	getsub := handlers.NewMessage(
		admin(message.HasPrefix("!get")),
		_getsub,
	)

	d.AddHandler(addsub)
	d.AddHandler(getsub)
}

func admin(cond func(*gotgbot.Message) bool) func(*gotgbot.Message) bool {
	return func(msg *gotgbot.Message) bool {
		return tools.IsAdmin(msg.From.Id) && cond(msg)
	}
}
