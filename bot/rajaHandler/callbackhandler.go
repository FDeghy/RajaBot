package rajaHandler

import (
	"RajaBot/core"
	"RajaBot/database"
	"RajaBot/tools"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	ptime "github.com/yaa110/go-persian-calendar"
)

// new
func _newcb(b *gotgbot.Bot, ctx *ext.Context) error {
	if mode := strings.TrimPrefix(ctx.CallbackQuery.Data, "new-"); mode == "raja" {
		markup, err := createStationsMarkup(0, "src-raja")
		if err != nil {
			b.SendMessage(ctx.EffectiveChat.Id, StationsLoadErr, nil)
			return nil
		}

		b.EditMessageText(SrcMsg, &gotgbot.EditMessageTextOpts{
			ChatId:      ctx.EffectiveChat.Id,
			MessageId:   ctx.EffectiveMessage.MessageId,
			ReplyMarkup: *markup,
		})
	} else if mode == "homei" {
		markup := createRoutesMarkup()

		b.EditMessageText(SelectRouteMSg, &gotgbot.EditMessageTextOpts{
			ChatId:      ctx.EffectiveChat.Id,
			MessageId:   ctx.EffectiveMessage.MessageId,
			ReplyMarkup: *markup,
		})
	} else if mode == "thrdapp" {
		markup, err := createStationsMarkup(0, "src-ta")
		if err != nil {
			b.SendMessage(ctx.EffectiveChat.Id, StationsLoadErr, nil)
			return nil
		}

		b.EditMessageText(SrcMsg, &gotgbot.EditMessageTextOpts{
			ChatId:      ctx.EffectiveChat.Id,
			MessageId:   ctx.EffectiveMessage.MessageId,
			ReplyMarkup: *markup,
		})
	}

	return nil
}

// stations
func _pgCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	var prefix string

	d := strings.TrimPrefix(ctx.CallbackQuery.Data, "pg-")
	if strings.HasPrefix(d, "src-raja") {
		prefix = "src-raja"
	} else if strings.HasPrefix(d, "src-ta") {
		prefix = "src-ta"
	} else if strings.HasPrefix(d, "dst-raja") {
		prefix = "dst-raja"
	} else if strings.HasPrefix(d, "dst-ta") {
		prefix = "dst-ta"
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

	if user.State != "normal" {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}
	if tools.CheckReachLimit(*user) {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, LimitReached, nil)
		return nil
	}
	if !tools.CheckHaveSubscription(user.UserID) {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, NoHaveSub, nil)
		return nil
	}

	if strings.HasPrefix(ctx.CallbackQuery.Data, "src-raja") {
		src, err := strconv.Atoi(strings.TrimPrefix(ctx.CallbackQuery.Data, "src-raja-"))
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}
		train := &database.TrainWR{
			UserID: user.UserID,
			Src:    src,
			IsDone: false,
		}
		tid := database.SaveTrainWR(train)

		user.State = fmt.Sprintf("src-raja-%d", tid)
		database.UpdateTgUser(user)

		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

		markup, err := createStationsMarkup(0, "dst-raja")
		if err != nil {
			b.SendMessage(ctx.EffectiveChat.Id, StationsLoadErr, nil)
			return nil
		}

		b.SendMessage(ctx.EffectiveChat.Id, DstMsg, &gotgbot.SendMessageOpts{
			ReplyMarkup: markup,
		})

	} else if strings.HasPrefix(ctx.CallbackQuery.Data, "src-ta") {
		src, err := strconv.Atoi(strings.TrimPrefix(ctx.CallbackQuery.Data, "src-ta-"))
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}
		train := &database.TrainWR{
			UserID:  user.UserID,
			Src:     src,
			IsDone:  false,
			ThrdApp: 1,
		}
		tid := database.SaveTrainWR(train)

		user.State = fmt.Sprintf("src-ta-%d", tid)
		database.UpdateTgUser(user)

		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

		markup, err := createStationsMarkup(0, "dst-ta")
		if err != nil {
			b.SendMessage(ctx.EffectiveChat.Id, StationsLoadErr, nil)
			return nil
		}

		b.SendMessage(ctx.EffectiveChat.Id, DstMsg, &gotgbot.SendMessageOpts{
			ReplyMarkup: markup,
		})

	}
	return nil
}

// destination select
func _dstCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	var dst int
	var tid uint64
	var train *database.TrainWR

	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	if !strings.HasPrefix(user.State, "src-") {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}

	if strings.HasPrefix(ctx.CallbackQuery.Data, "dst-raja") {
		_tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "src-raja-"), 10, 64)
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}
		tid = _tid
		train = database.GetTrainWRByTid(tid)
		if train == nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}

		dst, err = strconv.Atoi(strings.TrimPrefix(ctx.CallbackQuery.Data, "dst-raja-"))
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}

		train.Dst = dst
		database.UpdateTrainWR(train)
		user.State = fmt.Sprintf("dst-raja-%d", tid)
		database.UpdateTgUser(user)

	} else if strings.HasPrefix(ctx.CallbackQuery.Data, "dst-ta") {
		_tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "src-ta-"), 10, 64)
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}
		tid = _tid
		train = database.GetTrainWRByTid(tid)
		if train == nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}

		dst, err = strconv.Atoi(strings.TrimPrefix(ctx.CallbackQuery.Data, "dst-ta-"))
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}

		train.Dst = dst
		database.UpdateTrainWR(train)
		user.State = fmt.Sprintf("dst-ta-%d", tid)
		database.UpdateTgUser(user)

	}

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
	if strings.HasPrefix(user.State, "dst-raja-") {
		tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "dst-raja-"), 10, 64)
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
		user.State = fmt.Sprintf("taq-raja-%d", tid)
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

	} else if strings.HasPrefix(user.State, "dst-ta-") {
		tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "dst-ta-"), 10, 64)
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
		user.State = fmt.Sprintf("taq-ta-%d", tid)
		database.UpdateTgUser(user)

		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

		msg, _ := b.SendMessage(ctx.EffectiveChat.Id, GetTrainsInfoMsg, &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
		})
		// todo
		markup, err := createTrainListThrdAppMarkup(*train)
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

	} else if strings.HasPrefix(user.State, "rt-") {
		tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "rt-"), 10, 64)
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
		user.State = fmt.Sprintf("rttaq-%d", tid)
		database.UpdateTgUser(user)

		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

		msg, _ := b.SendMessage(ctx.EffectiveChat.Id, GetTrainsInfoMsg, &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{},
		})

		markup, err := createTrainRtListMarkup(*train)
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

	} else {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
	}
	return nil
}

// train & time (new raja site) select
func _rttrCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	if !strings.HasPrefix(user.State, "rttaq-") {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}

	// get train record by user state
	tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "rttaq-"), 10, 64)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	train := database.GetTrainWRByTid(tid)
	if train == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	// get train id and startTime from callback data
	d := strings.TrimPrefix(ctx.CallbackQuery.Data, "rttr-")
	data := strings.Split(d, "-")
	trainId, err := strconv.Atoi(data[0])
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	startTime := data[1]

	b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

	// update user state to normal and check subscription
	user.State = "normal"
	database.UpdateTgUser(user)
	if !tools.CheckHaveSubscription(user.UserID) {
		b.SendMessage(ctx.EffectiveChat.Id, NoHaveSub, nil)
		return nil
	}

	//send to core
	err = core.HandleGoFetch(train)
	if err != nil {
		b.SendMessage(ctx.EffectiveChat.Id, RajaErr, nil)
		return nil
	}

	train.TrainId = trainId
	train.Hour = startTime
	database.UpdateTrainWR(train)
	b.SendMessage(ctx.EffectiveChat.Id, successfulCreate, nil)

	return nil
}

// train (& time) select
func _trCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	var trainId int
	var tid uint64
	var train *database.TrainWR
	var data []string

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

	if strings.HasPrefix(user.State, "taq-raja-") {
		_tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "taq-raja-"), 10, 64)
		tid = _tid
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}
		_train := database.GetTrainWRByTid(tid)
		train = _train
		if train == nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}

		d := strings.TrimPrefix(ctx.CallbackQuery.Data, "tr-raja-")
		_data := strings.Split(d, "-")
		data = _data
		_trainId, err := strconv.Atoi(data[0])
		trainId = _trainId
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}

	} else if strings.HasPrefix(user.State, "taq-ta-") {
		_tid, err := strconv.ParseUint(strings.TrimPrefix(user.State, "taq-ta-"), 10, 64)
		tid = _tid
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}
		_train := database.GetTrainWRByTid(tid)
		train = _train
		if train == nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}

		d := strings.TrimPrefix(ctx.CallbackQuery.Data, "tr-ta-")
		_data := strings.Split(d, "-")
		data = _data
		_trainId, err := strconv.Atoi(data[0])
		trainId = _trainId
		if err != nil {
			b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
			return nil
		}

	}
	user.State = "normal"
	database.UpdateTgUser(user)

	b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

	if !tools.CheckHaveSubscription(user.UserID) {
		b.SendMessage(ctx.EffectiveChat.Id, NoHaveSub, nil)
		return nil
	}

	//send to core
	err := core.HandleGoFetch(train)
	if err != nil {
		b.SendMessage(ctx.EffectiveChat.Id, RajaErr, nil)
		return nil
	}

	train.TrainId = trainId
	train.Hour = data[1]
	database.UpdateTrainWR(train)

	var inlineKey = [][]gotgbot.InlineKeyboardButton{{
		{
			Text: core.RajaSearchButTxt,
			Url: fmt.Sprintf(
				core.RajaSearchURL,
				train.Src,
				train.Dst,
				ptime.Unix(train.Day, 0).Format(core.RajaSearchDateFmt),
			),
		},
	}}
	if train.ThrdApp != 0 {
		inlineKey = nil
	}

	b.SendMessage(ctx.EffectiveChat.Id, successfulCreate, &gotgbot.SendMessageOpts{
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: inlineKey,
		},
	})

	return nil
}

// old train (& time) select
func _oldtrCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: OldTrErr, ShowAlert: true})
	return nil
}

func _cancCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	d := strings.TrimPrefix(ctx.CallbackQuery.Data, "canc-")
	id, err := strconv.ParseUint(d, 10, 64)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	err = core.CancelWork(id)
	if err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: err.Error(), ShowAlert: true})
		return nil
	}
	b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: CancOkAlert, ShowAlert: true})

	trainWRs := database.GetActiveTrainWRs(user.UserID)
	markup := createListMarkup(trainWRs)
	b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
	b.SendMessage(ctx.EffectiveChat.Id, createListMsg(trainWRs), &gotgbot.SendMessageOpts{
		ParseMode:   gotgbot.ParseModeMarkdownV2,
		ReplyMarkup: markup,
	})

	return nil
}

// nil
func _nilCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: NilButton, ShowAlert: true})
	return nil
}

// route select
func _rtCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	user := database.GetTgUser(ctx.EffectiveSender.Id())
	if user == nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}

	if user.State != "normal" {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, StateErr, nil)
		return nil
	}
	if tools.CheckReachLimit(*user) {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, LimitReached, nil)
		return nil
	}
	if !tools.CheckHaveSubscription(user.UserID) {
		b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)
		b.SendMessage(ctx.EffectiveChat.Id, NoHaveSub, nil)
		return nil
	}

	rtId := strings.TrimPrefix(ctx.CallbackQuery.Data, "rt-")
	rt := Routes.FindRoute(rtId)
	src, err := strconv.Atoi(rtId)
	if rt == nil || err != nil {
		b.AnswerCallbackQuery(ctx.CallbackQuery.Id, &gotgbot.AnswerCallbackQueryOpts{Text: AnError, ShowAlert: true})
		return nil
	}
	train := &database.TrainWR{
		UserID: user.UserID,
		Src:    src,
		Dst:    -1,
		IsDone: false,
	}
	tid := database.SaveTrainWR(train)

	user.State = fmt.Sprintf("rt-%d", tid)
	database.UpdateTgUser(user)

	b.DeleteMessage(ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId, nil)

	now := ptime.Now()
	markup, _ := createTaqvimMarkup(now.Year(), int(now.Month()))
	b.SendMessage(ctx.EffectiveChat.Id, now.Format(DayMsg), &gotgbot.SendMessageOpts{
		ReplyMarkup: markup,
	})
	return nil
}
