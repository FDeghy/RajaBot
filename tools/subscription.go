package tools

import (
	"RajaBot/config"
	"RajaBot/database"
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	ptime "github.com/yaa110/go-persian-calendar"
)

func CheckHaveSubscription(userId int64) bool {
	sub := database.GetSubscription(userId)
	if sub == nil {
		return false
	}
	if IsAdmin(userId) || (sub.IsEnabled && sub.ExpirationDate > time.Now().Unix()) {
		return true
	}
	return false
}

func CreateSubStatus(sub database.Subscription) (string, *gotgbot.InlineKeyboardMarkup) {
	markup := &gotgbot.InlineKeyboardMarkup{}
	if sub.IsTrial {
		markup.InlineKeyboard = [][]gotgbot.InlineKeyboardButton{
			{
				{Text: BuySub, CallbackData: "buysub"},
			},
		}
	} else {
		markup.InlineKeyboard = [][]gotgbot.InlineKeyboardButton{
			{
				{Text: BuySub, CallbackData: "buysub"},
				{Text: FreeTrial, CallbackData: "freetrial"},
			},
		}
	}

	status := Disabled
	if sub.IsEnabled {
		status = Enabled
	}
	registeryDate := Unkown
	if sub.RegisteryDate != 0 {
		registeryDate = ptime.Unix(sub.RegisteryDate, 0).Format(TimeFormat)
	}
	expirationDate := Disabled
	if sub.ExpirationDate != 0 {
		expirationDate = ptime.Unix(sub.RegisteryDate, 0).Format(TimeFormat)
	}

	return fmt.Sprintf(
		SubscriptionStatusMsg,
		status,
		registeryDate,
		expirationDate,
		config.Cfg.Bot.UserLimit,
	), markup
}
