package tools

import (
	"RajaBot/config"
	"RajaBot/database"
	"fmt"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func AppendEmptyButton(r *[]gotgbot.InlineKeyboardButton, n int) {
	for i := 0; i < n; i++ {
		*r = append(*r, gotgbot.InlineKeyboardButton{
			Text:         " ",
			CallbackData: "nil",
		})
	}
}

// 100000 -> 100,000
func NumToMoney(num int) string {
	n := []byte(fmt.Sprint(num))
	slices.Reverse(n)
	var res []byte
	for i, r := range n {
		if i%3 == 0 && i != 0 {
			res = append(res, ',')
		}
		res = append(res, r)
	}
	slices.Reverse(res)
	return string(res)
}

func CheckReachLimit(user database.TgUser) bool {
	limit := config.Cfg.Bot.UserLimit
	if user.IsVip {
		limit = config.Cfg.Bot.VipLimit
	}
	activeTrains := database.GetActiveTrainWRs(user.UserID)
	return len(activeTrains) >= limit
}

func IsAdmin(userId int64) bool {
	return slices.Contains(config.Cfg.Bot.Admins, userId)
}
