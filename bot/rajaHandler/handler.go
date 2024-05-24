package rajaHandler

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
)

func Load(d *ext.Dispatcher) {
	start := handlers.NewCommand("start", _start)
	new := handlers.NewCommand("new", _new)
	cancel := handlers.NewCommand("cancel", _cancel)
	list := handlers.NewCommand("list", _list)
	test := handlers.NewCommand("test", _test)
	newCallback := handlers.NewCallback(callbackquery.Prefix("new-"), _newcb)
	pgCallback := handlers.NewCallback(callbackquery.Prefix("pg-"), _pgCallback)
	pgmCallback := handlers.NewCallback(callbackquery.Prefix("pgm-"), _pgmCallback)
	srcCallback := handlers.NewCallback(callbackquery.Prefix("src-"), _srcCallback)
	dstCallback := handlers.NewCallback(callbackquery.Prefix("dst-"), _dstCallback)
	taqCallback := handlers.NewCallback(callbackquery.Prefix("taq-"), _taqCallback)
	trCallback := handlers.NewCallback(callbackquery.Prefix("tr-"), _trCallback)
	oldtrCallback := handlers.NewCallback(callbackquery.Prefix("oldtr-"), _oldtrCallback)
	cancCallback := handlers.NewCallback(callbackquery.Prefix("canc-"), _cancCallback)
	nilCallback := handlers.NewCallback(callbackquery.Equal("nil"), _nilCallback)
	rtCallback := handlers.NewCallback(callbackquery.Prefix("rt-"), _rtCallback)
	rttrCallback := handlers.NewCallback(callbackquery.Prefix("rttr-"), _rttrCallback)

	d.AddHandler(start)
	d.AddHandler(new)
	d.AddHandler(cancel)
	d.AddHandler(list)
	d.AddHandler(test)
	d.AddHandler(newCallback)
	d.AddHandler(pgCallback)
	d.AddHandler(pgmCallback)
	d.AddHandler(srcCallback)
	d.AddHandler(dstCallback)
	d.AddHandler(taqCallback)
	d.AddHandler(trCallback)
	d.AddHandler(oldtrCallback)
	d.AddHandler(cancCallback)
	d.AddHandler(nilCallback)
	d.AddHandler(rtCallback)
	d.AddHandler(rttrCallback)
}
