package bot

import (
	"RajaBot/config"
	"RajaBot/database"
	"fmt"
	"slices"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func favoriteSign(st string) string {
	for _, i := range config.Cfg.Bot.FavoriteStations {
		if strings.EqualFold(st, i) {
			return FavSign
		}
	}
	return ""
}

func appendEmptyButton(r *[]gotgbot.InlineKeyboardButton, n int) {
	for i := 0; i < n; i++ {
		*r = append(*r, gotgbot.InlineKeyboardButton{
			Text:         " ",
			CallbackData: "nil",
		})
	}
}

// 100000 -> 100,000
func numToMoney(num int) string {
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

func checkLimit(user database.TgUser) bool {
	limit := config.Cfg.Bot.UserLimit
	if user.IsVip {
		limit = config.Cfg.Bot.VipLimit
	}
	activeTrains := database.GetActiveTrainWRs(user.UserID)
	return len(*activeTrains) >= limit
}

func escapeMarkdown(inp string) string {
	chars := []rune{'_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'}
	result := inp
	for _, c := range chars {
		result = strings.ReplaceAll(result, string(c), "\\"+string(c))
	}
	return result
}
