package rajaHandler

import (
	"RajaBot/database"
	"RajaBot/tools"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func _start(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		user = database.NewTgUser(ctx.EffectiveSender.Id())
		database.SaveTgUser(user)
	}

	sub := database.GetSubscription(ctx.EffectiveSender.Id())
	if sub == nil {
		sub = database.NewSubscription(user.UserID)
		database.SaveSubscription(sub)
	}

	b.SendMessage(ctx.EffectiveChat.Id, StartMsg, nil)
	return nil
}

func _new(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.SendMessage(ctx.EffectiveChat.Id, AnError, nil)
		return nil
	}

	if user.State != "normal" {
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}
	if tools.CheckReachLimit(*user) {
		b.SendMessage(ctx.EffectiveChat.Id, LimitReached, nil)
		return nil
	}
	if !tools.CheckHaveSubscription(user.UserID) {
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

	if Stations == nil {//|| Routes == nil {
		b.SendMessage(ctx.EffectiveChat.Id, StationsLoadErr, nil)
		return nil
	}

	b.SendMessage(ctx.EffectiveChat.Id, SelectTypeMsg, &gotgbot.SendMessageOpts{
		ReplyParameters: &gotgbot.ReplyParameters{
			MessageId: ctx.EffectiveMessage.MessageId,
		},
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: NormalRajaBtn, CallbackData: "new-raja"}}, //{Text: HomeiBtn, CallbackData: "new-homei"}},
				{{Text: ThrdAppBtn, CallbackData: "new-thrdapp"}},
			},
		},
	})

	return nil
}

func _cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.SendMessage(ctx.EffectiveChat.Id, AnError, nil)
		return nil
	}

	if user.State == "normal" {
		b.SendMessage(ctx.EffectiveChat.Id, CancelOkMsg, nil)
		return nil
	}
	d := strings.Split(user.State, "-")
	tid, err := strconv.ParseUint(d[len(d)-1], 10, 64)
	if err != nil {
		b.SendMessage(ctx.EffectiveChat.Id, AnError, nil)
		return nil
	}
	train := database.GetTrainWRByTid(tid)
	if train == nil {
		b.SendMessage(ctx.EffectiveChat.Id, AnError, nil)
		return nil
	}
	database.DeleteTrainWR(train)

	user.State = "normal"
	database.UpdateTgUser(user)

	b.SendMessage(ctx.EffectiveChat.Id, CancelOkMsg, nil)

	return nil
}

func _list(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.SendMessage(ctx.EffectiveChat.Id, AnError, nil)
		return nil
	}

	if user.State != "normal" {
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}

	trainWRs := database.GetActiveTrainWRs(user.UserID)
	markup := createListMarkup(trainWRs)
	b.SendMessage(ctx.EffectiveChat.Id, createListMsg(trainWRs), &gotgbot.SendMessageOpts{
		ParseMode:   gotgbot.ParseModeMarkdownV2,
		ReplyMarkup: markup,
	})
	return nil
}

func _test(b *gotgbot.Bot, ctx *ext.Context) error {
	markup, _ := createTaqvimMarkup(1403, 1)
	b.SendMessage(ctx.EffectiveChat.Id, "test", &gotgbot.SendMessageOpts{
		ReplyMarkup: markup,
	})
	return nil
}
