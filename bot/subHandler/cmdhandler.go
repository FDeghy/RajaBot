package subHandler

import (
	"RajaBot/database"
	"RajaBot/tools"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func _sub(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.SendMessage(ctx.EffectiveChat.Id, AnError, nil)
		return nil
	}

	sub := database.GetSubscription(ctx.EffectiveUser.Id)
	if sub == nil {
		sub = database.NewSubscription(user.UserID)
		database.SaveSubscription(sub)
	}

	text, markup := tools.CreateSubStatus(sub)

	b.SendMessage(
		ctx.EffectiveChat.Id,
		text,
		&gotgbot.SendMessageOpts{
			ReplyParameters: &gotgbot.ReplyParameters{
				MessageId: ctx.EffectiveMessage.MessageId,
			},
			ReplyMarkup: markup,
		},
	)

	return nil
}
