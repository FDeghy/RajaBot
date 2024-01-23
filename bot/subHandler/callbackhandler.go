package subHandler

import (
	"RajaBot/config"
	"RajaBot/database"
	"RajaBot/tools"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	ptime "github.com/yaa110/go-persian-calendar"
)

func _buysub(b *gotgbot.Bot, ctx *ext.Context) error {
	return nil
}

func _freetrial(b *gotgbot.Bot, ctx *ext.Context) error {
	sub := database.GetSubscription(ctx.EffectiveUser.Id)
	if sub == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	if sub.IsTrial {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AlreadyTrial, ShowAlert: true})
		return nil
	}

	if sub.ExpirationDate == 0 {
		sub.ExpirationDate = time.Now().Unix()
	}
	sub.ExpirationDate = ptime.Unix(sub.ExpirationDate, 0).AddDate(0, 0, config.Cfg.Bot.TrialDays).Unix()
	sub.IsTrial = true
	sub.IsEnabled = true
	database.UpdateSubscription(sub)

	b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: EnabledFreeTrial, ShowAlert: true})
	text, markup := tools.CreateSubStatus(*sub)

	b.EditMessageText(
		text,
		&gotgbot.EditMessageTextOpts{
			ChatId:      ctx.EffectiveChat.Id,
			MessageId:   ctx.EffectiveMessage.MessageId,
			ReplyMarkup: *markup,
		},
	)
	return nil
}
