package bot

import (
	"RajaBot/bot/rajaHandler"
	"RajaBot/config"
	"RajaBot/core"
	"RajaBot/payment"
	"RajaBot/tools"
	"RajaBot/tools/tlog"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func CreateBot() error {
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

	b, err := gotgbot.NewBot(config.Cfg.Bot.Token, &gotgbot.BotOpts{
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
	bot = b
	// load stations
	err = loadStations()
	if err != nil {
		return errors.New("failed to load stations")
	}
	// load routes
	//err = loadRoutes()
	//if err != nil {
	//	return errors.New("failed to load routes")
	//}

	dispatcher = ext.NewDispatcher(&ext.DispatcherOpts{})
	updater = ext.NewUpdater(dispatcher, nil)

	// set Bot in core
	core.Bot = bot
	// set Bot in payment
	payment.Bot = bot
	// set Bot in tools
	tools.Bot = bot
	// set Bot in tlog
	tlog.Bot = bot

	return nil
}

func StartBot() error {
	err := updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
	})
	if err != nil {
		return errors.New("failed to start polling updates")
	}

	log.Printf("Bot -> Bot started. %v - %v\n", bot.Username, bot.Id)

	// load handlers
	load(dispatcher)

	updater.Idle()
	return nil
}

func loadStations() error {
	sts, err := raja.GetStations()
	if err != nil {
		return err
	}
	rajaHandler.Stations = sts
	return nil
}

func loadRoutes() error {
	err := tools.LoadRoutes()
	rajaHandler.Routes = tools.Routes
	return err
}
