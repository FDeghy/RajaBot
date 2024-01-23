package subHandler

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

func Load(d *ext.Dispatcher) {
	sub := handlers.NewCommand("sub", _sub)
	buysub := handlers.NewCallback(callbackquery.Equal("buysub"), _buysub)
	freetrial := handlers.NewCallback(callbackquery.Equal("freetrial"), _freetrial)

	d.AddHandler(sub)
	d.AddHandler(buysub)
	d.AddHandler(freetrial)
}