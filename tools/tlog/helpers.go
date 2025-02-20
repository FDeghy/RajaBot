package tlog

import (
	"RajaBot/config"
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	ptime "github.com/yaa110/go-persian-calendar"
)

func userMarkdown(userId int64) (string, error) {
	chat, err := Bot.GetChat(userId, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"<a href=\"tg://openmessage?user_id=%v\">(%v %v)</a> (%v) @%v",
		chat.Id,
		chat.FirstName,
		chat.LastName,
		chat.Id,
		chat.Username,
	), nil
}

func SendLog(userId int64, mode Mode, opt string) {
	userMD, err := userMarkdown(userId)
	if err != nil {
		log.Println("cannot GetChat in SengLog")
		return
	}

	if mode == Payment {
		Bot.SendMessage(
			config.Cfg.Bot.LogChannel,
			fmt.Sprintf(PaymentMsg, userMD, opt, ptime.Now().Format(DateFmt)),
			&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML},
		)
	} else if mode == NewTrain {
		Bot.SendMessage(
			config.Cfg.Bot.LogChannel,
			fmt.Sprintf(NewTrainMsg, userMD, opt, ptime.Now().Format(DateFmt)),
			&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML},
		)
	}
}
