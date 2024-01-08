package bot

import (
	"RajaBot/config"
	"RajaBot/core"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func StartBot() error {
	httpClient := &http.Client{}
	if config.Cfg.Bot.HttpURI != "" {
		proxy, err := url.Parse(config.Cfg.Bot.HttpURI)
		if err != nil {
			return errors.New("failed to parse http proxy")
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
	}

	bot, err := gotgbot.NewBot(config.Cfg.Bot.Token, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: *httpClient,
		},
		RequestOpts: &gotgbot.RequestOpts{
			Timeout: time.Duration(config.Cfg.Bot.Timeout) * gotgbot.DefaultTimeout,
		},
	})
	if err != nil {
		return errors.New("failed to create bot client")
	}

	// load stations
	err = loadStations()
	if err != nil {
		return errors.New("failed to load stations")
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{})
	updater := ext.NewUpdater(dispatcher, nil)
	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
	})
	if err != nil {
		return errors.New("failed to start polling updates")
	}

	log.Printf("Bot started. %v - %v\n", bot.Username, bot.Id)

	// load handlers
	load(dispatcher)
	// set Bot in core
	core.Bot = bot

	updater.Idle()
	return nil
}

func loadStations() error {
	sts, err := raja.GetStations()
	if err != nil {
		return err
	}
	Stations = sts
	return nil
}
