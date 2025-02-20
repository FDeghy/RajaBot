package tools

import (
	"RajaBot/config"
	"RajaBot/database"
	"RajaBot/tools/tlog"
	"fmt"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	ptime "github.com/yaa110/go-persian-calendar"
)

func CheckHaveSubscription(userId int64) bool {
	sub := database.GetSubscription(userId)
	if sub == nil {
		return false
	}
	if IsAdmin(userId) || sub.ExpirationDate > time.Now().Unix() {
		sub.IsEnabled = true
	} else {
		sub.IsEnabled = false
	}
	database.UpdateSubscription(sub)
	return sub.IsEnabled
}

func CreateSubStatus(sub *database.Subscription) (string, *gotgbot.InlineKeyboardMarkup) {
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
			},
			{
				{Text: FreeTrial, CallbackData: "freetrial"},
			},
		}
	}

	sub.IsEnabled = false
	status := Disabled
	if sub.ExpirationDate > time.Now().Unix() {
		sub.IsEnabled = true
		status = Enabled
	}
	database.UpdateSubscription(sub)
	registeryDate := Unkown
	if sub.RegisteryDate != 0 {
		registeryDate = ptime.Unix(sub.RegisteryDate, 0).Format(TimeFormat)
	}
	expirationDate := Disabled
	if sub.ExpirationDate != 0 {
		expirationDate = ptime.Unix(sub.ExpirationDate, 0).Format(TimeFormat)
	}

	return fmt.Sprintf(
		SubscriptionStatusMsg,
		status,
		registeryDate,
		expirationDate,
		config.Cfg.Bot.UserLimit,
	), markup
}

func AddDaysSub(userId int64, days int) (*database.Subscription, error) {
	sub := database.GetSubscription(userId)
	if sub == nil {
		return sub, ErrSubNotFound
	}
	if sub.ExpirationDate < time.Now().Unix() {
		sub.ExpirationDate = time.Now().Unix()
	}
	sub.ExpirationDate = ptime.Unix(sub.ExpirationDate, 0).AddDate(0, 0, days).Unix()
	sub.IsEnabled = true
	database.UpdateSubscription(sub)
	tlog.SendLog(userId, tlog.Payment, strconv.Itoa(days)+" روز")
	return sub, nil
}

func SetTrialSub(userId int64) (*database.Subscription, error) {
	sub := database.GetSubscription(userId)
	if sub == nil {
		return sub, ErrSubNotFound
	}
	if sub.IsTrial {
		return sub, ErrAlreadyTrial
	}
	if sub.ExpirationDate < time.Now().Unix() {
		sub.ExpirationDate = time.Now().Unix()
	}
	sub.ExpirationDate = ptime.Unix(sub.ExpirationDate, 0).AddDate(0, 0, config.Cfg.Bot.TrialDays).Unix()
	sub.IsTrial = true
	sub.IsEnabled = true
	database.UpdateSubscription(sub)
	tlog.SendLog(userId, tlog.Payment, "فیری تریال")
	return sub, nil
}
