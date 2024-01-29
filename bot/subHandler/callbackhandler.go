package subHandler

import (
	"RajaBot/config"
	"RajaBot/database"
	"RajaBot/payment"
	"RajaBot/tools"
	"errors"
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func _freetrial(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	sub, err := tools.SetTrialSub(user.UserID)
	if errors.Is(err, tools.ErrSubNotFound) {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	if errors.Is(err, tools.ErrAlreadyTrial) {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AlreadyTrial, ShowAlert: true})
		return nil
	}

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

func _buysub(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	sub := database.GetSubscription(user.UserID)
	if sub == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	paym, err := payment.NewTransaction(user, config.Cfg.Payment.OneMonthPrice)
	if err != nil {
		if errors.Is(err, payment.ErrUncompletedTransactionFound) {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: UncompletedTransaction, ShowAlert: true})
		} else if errors.Is(err, payment.ErrBadCode) {
			log.Printf("Bot -> new transaction error: %v", err)
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		}
		return nil
	}
	b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: TransactionCreated, ShowAlert: true})
	b.SendMessage(
		ctx.EffectiveChat.Id,
		GoTransaction,
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
					{
						Text: fmt.Sprintf(OneMonthBtn, tools.NumToMoney(int(config.Cfg.Payment.OneMonthPrice))),
						Url:  payment.CreateBankLink(paym.TransID),
					},
				}},
			},
		},
	)
	return nil
}
