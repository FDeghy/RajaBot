package bot

import (
	"RajaBot/bot/adminHandler"
	"RajaBot/bot/rajaHandler"
	"RajaBot/bot/subHandler"
	"RajaBot/config"

	"github.com/ALiwoto/ratelimiter"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func load(d *ext.Dispatcher) {
	loadLimiter(d)
	rajaHandler.Load(d)
	subHandler.Load(d)
	adminHandler.Load(d)
}

func loadLimiter(d *ext.Dispatcher) {
	rateLimiter := ratelimiter.NewLimiter(d, nil)
	rateLimiter.TextOnly = true
	rateLimiter.AddExceptionID(config.Cfg.Bot.Admins...)
	rateLimiter.Start()
}
