package bot

import (
	"RajaBot/core"
	"RajaBot/database"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	ptime "github.com/yaa110/go-persian-calendar"
)

// stations
func _pgCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	var prefix string

	d := strings.TrimPrefix(ctx.CallbackQuery.Data, "pg-")
	if strings.HasPrefix(d, "src") {
		prefix = "src"
	} else if strings.HasPrefix(d, "dst") {
		prefix = "dst"
	} else {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	d = strings.TrimPrefix(d, prefix+"-")
	page, err := strconv.Atoi(d)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	markup, err := createStationsMarkup(page, prefix)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	b.EditMessageReplyMarkup(&gotgbot.EditMessageReplyMarkupOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: *markup,
	})
	b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{
		Text: fmt.Sprintf(PageN, page),
	})
	return nil
}

// calender
func _pgmCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	d := strings.TrimPrefix(ctx.CallbackQuery.Data, "pgm-")
	unix, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	pt := ptime.Unix(unix, 0)
	markup, _ := createTaqvimMarkup(pt.Year(), int(pt.Month()))
	b.EditMessageText(pt.Format(DayMsg), &gotgbot.EditMessageTextOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   ctx.EffectiveMessage.MessageId,
		ReplyMarkup: *markup,
	})
	return nil
}

// source select
func _srcCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	if checkLimit(*user) {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, LimitReached, nil)
		return nil
	}

	src, err := strconv.Atoi(strings.TrimPrefix(ctx.CallbackQuery.Data, "src-"))
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	train := &database.TrainWR{
		UserID: user.UserID,
		Src:    src,
		IsDone: false,
	}
	tid, err := database.SaveTrainWR(train)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	user.State = fmt.Sprintf("src-%d", tid)
	database.UpdateTgUser(user)

	b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

	markup, err := createStationsMarkup(0, "dst")
	if err != nil {
		b.SendMessage(ctx.EffectiveChat.Id, StationsLoadErr, nil)
		return nil
	}

	b.SendMessage(ctx.EffectiveChat.Id, DstMsg, &gotgbot.SendMessageOpts{
		ReplyMarkup: markup,
	})
	return nil
}

// destination select
func _dstCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	if !strings.HasPrefix(user.State, "src-") {
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}

	tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "src-"), 10, 64)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	train := database.GetTrainWRByTid(tid)
	if train == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	dst, err := strconv.Atoi(strings.TrimPrefix(ctx.CallbackQuery.Data, "dst-"))
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	train.Dst = dst
	database.UpdateTrainWR(train)
	user.State = fmt.Sprintf("dst-%d", tid)
	database.UpdateTgUser(user)

	b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

	now := ptime.Now()
	markup, _ := createTaqvimMarkup(now.Year(), int(now.Month()))
	b.SendMessage(ctx.EffectiveChat.Id, now.Format(DayMsg), &gotgbot.SendMessageOpts{
		ReplyMarkup: markup,
	})

	return nil
}

// date select
func _taqCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	if !strings.HasPrefix(user.State, "dst-") {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}

	tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "dst-"), 10, 64)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	train := database.GetTrainWRByTid(tid)
	if train == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	d := strings.TrimPrefix(ctx.CallbackQuery.Data, "taq-")
	unix, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	lastSecDay := ptime.Unix(unix, 0)
	lastSecDay.Set(lastSecDay.Year(), lastSecDay.Month(), lastSecDay.Day(), 23, 59, 59, 0, ptime.Iran())
	if time.Now().Unix() >= lastSecDay.Unix() {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: OldDateErr, ShowAlert: true})
		return nil
	}

	pt := ptime.Unix(unix, 0)
	train.Day = unix
	database.UpdateTrainWR(train)
	user.State = fmt.Sprintf("taq-%d", tid)
	database.UpdateTgUser(user)

	b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

	msg, _ := b.SendMessage(ctx.EffectiveChat.Id, GetTrainsInfoMsg, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
	})
	markup, err := createTrainListMarkup(*train)
	if err != nil {
		b.EditMessageText(TrainListLoadErr, &gotgbot.EditMessageTextOpts{
			ChatId:      ctx.EffectiveChat.Id,
			MessageId:   msg.MessageId,
			ReplyMarkup: *markup,
		})
		b.SendMessage(ctx.EffectiveChat.Id, CancelMsg, nil)
		return nil
	}
	b.EditMessageText(pt.Format(TrainSelMsg), &gotgbot.EditMessageTextOpts{
		ChatId:      ctx.EffectiveChat.Id,
		MessageId:   msg.MessageId,
		ReplyMarkup: *markup,
	})
	return nil
}

// train (& time) select
func _trCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	if !strings.HasPrefix(user.State, "taq-") {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}

	tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "taq-"), 10, 64)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	train := database.GetTrainWRByTid(tid)
	if train == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	d := strings.TrimPrefix(ctx.CallbackQuery.Data, "tr-")
	data := strings.Split(d, "-")
	trainId, err := strconv.Atoi(data[0])
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
	//send to core
	err = core.HandleGoFetch(train)
	if err != nil {
		b.SendMessage(ctx.EffectiveChat.Id, RajaErr, nil)
	} else {
		train.TrainId = trainId
		train.Hour = data[1]
		database.UpdateTrainWR(train)
		user.State = "normal"
		database.UpdateTgUser(user)
		b.SendMessage(ctx.EffectiveChat.Id, successfulCreate, nil)
	}

	return nil
}

// old train (& time) select
func _oldtrCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: OldTrErr, ShowAlert: true})
	return nil
}

// nil
func _nilCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: NilButton, ShowAlert: true})
	return nil
}
