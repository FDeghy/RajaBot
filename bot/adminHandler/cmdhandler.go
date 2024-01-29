package adminHandler

import (
	"RajaBot/database"
	"RajaBot/tools"
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	ptime "github.com/yaa110/go-persian-calendar"
)

func _addsub(b *gotgbot.Bot, ctx *ext.Context) error {
	data := strings.Split(ctx.EffectiveMessage.Text, " ")
	if len(data) != 3 {
		ctx.EffectiveMessage.Reply(b, "bad msg", nil)
		return ext.EndGroups
	}
	userId, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return ext.EndGroups
	}
	days, err := strconv.Atoi(data[2])
	if err != nil {
		ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return ext.EndGroups
	}

	_, err = tools.AddDaysSub(userId, days)
	if err != nil {
		ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return ext.EndGroups
	}
	b.SendMessage(userId, fmt.Sprintf(AddSubMsg, days), nil)
	ctx.EffectiveMessage.Reply(b, fmt.Sprintf("%v -> <a href=\"tg://user?id=%v\">%v</a>", days, userId, userId), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	return ext.EndGroups
}

func _getsub(b *gotgbot.Bot, ctx *ext.Context) error {
	data := strings.Split(ctx.EffectiveMessage.Text, " ")
	if len(data) != 2 {
		ctx.EffectiveMessage.Reply(b, "bad msg", nil)
		return ext.EndGroups
	}
	userId, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		ctx.EffectiveMessage.Reply(b, err.Error(), nil)
		return ext.EndGroups
	}

	sub := database.GetSubscription(userId)
	if sub == nil {
		ctx.EffectiveMessage.Reply(b, "not found", nil)
		return ext.EndGroups
	}
	ctx.EffectiveMessage.Reply(b, fmt.Sprintf("%#v", sub), nil)
	ctx.EffectiveMessage.Reply(b, ptime.Unix(sub.ExpirationDate, 0).Format("kk:mm E dd MMM yyyy"), nil)
	return ext.EndGroups
}
