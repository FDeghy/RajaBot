package bot

import (
	"RajaBot/config"
	"RajaBot/database"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
	ptime "github.com/yaa110/go-persian-calendar"
)

func createStationsMarkup(page int, prefix string) (*gotgbot.InlineKeyboardMarkup, error) {
	var sts raja.Stations
	markup := &gotgbot.InlineKeyboardMarkup{}
	row := []gotgbot.InlineKeyboardButton{}
	if (page+1)*config.Cfg.Bot.StationsButtonsPerPage >= len(*Stations) {
		sts = (*Stations)[page*config.Cfg.Bot.StationsButtonsPerPage:]
	} else {
		sts = (*Stations)[page*config.Cfg.Bot.StationsButtonsPerPage : (page+1)*config.Cfg.Bot.StationsButtonsPerPage]
	}
	for i, st := range sts {
		if i%4 == 0 && i != 0 {
			markup.InlineKeyboard = append(markup.InlineKeyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%s%s", st.PersianName, favoriteSign(st.EnglishName)),
			CallbackData: fmt.Sprintf("%v-%d", prefix, st.Id),
		})
		i++
	}
	if len(row) != 0 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}

	// page buttons
	if (page+1)*config.Cfg.Bot.StationsButtonsPerPage >= len(*Stations) {
		row = []gotgbot.InlineKeyboardButton{
			{Text: PreviousPage, CallbackData: fmt.Sprintf("pg-%s-%d", prefix, page-1)},
		}
	} else if page == 0 {
		row = []gotgbot.InlineKeyboardButton{
			{Text: NextPage, CallbackData: fmt.Sprintf("pg-%s-%d", prefix, page+1)},
		}
	} else {
		row = []gotgbot.InlineKeyboardButton{
			{Text: NextPage, CallbackData: fmt.Sprintf("pg-%s-%d", prefix, page+1)},
			{Text: PreviousPage, CallbackData: fmt.Sprintf("pg-%s-%d", prefix, page-1)},
		}
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	return markup, nil
}

func createTaqvimMarkup(sal int, mah int) (*gotgbot.InlineKeyboardMarkup, error) {
	markup := &gotgbot.InlineKeyboardMarkup{}
	weekDaysName := []gotgbot.InlineKeyboardButton{}
	for i := 0; i <= 6; i++ {
		weekDaysName = append(weekDaysName, gotgbot.InlineKeyboardButton{
			Text:         ptime.Weekday(i).Short(),
			CallbackData: "nil",
		})
	}
	slices.Reverse(weekDaysName)
	markup.InlineKeyboard = append(markup.InlineKeyboard, weekDaysName)
	curMah := ptime.Date(sal, ptime.Month(mah), 1, 0, 0, 0, 0, ptime.Iran())
	lastDayofMonth := curMah.LastMonthDay().Day()
	now := ptime.Now()
	pad := int(curMah.Weekday())
	row := []gotgbot.InlineKeyboardButton{}
	appendEmptyButton(&row, pad)
	for i := 1; i <= lastDayofMonth; i++ {
		day := gotgbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%d", i),
			CallbackData: fmt.Sprintf("taq-%d", curMah.Unix()),
		}
		if now.Year() == curMah.Year() && now.Month() == curMah.Month() && now.Day() == curMah.Day() {
			day.Text = fmt.Sprintf("(%d)", i)
		}
		row = append(row, day)
		if (pad+i)%7 == 0 {
			slices.Reverse(row)
			markup.InlineKeyboard = append(markup.InlineKeyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
		curMah = curMah.Tomorrow()
	}
	curMah = ptime.Date(sal, ptime.Month(mah), 1, 0, 0, 0, 0, ptime.Iran())
	pad = (7 - ((pad + lastDayofMonth) % 7)) % 7
	appendEmptyButton(&row, pad)
	if len(row) != 0 {
		slices.Reverse(row)
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}

	// page buttons
	if now.Month() == curMah.Month() {
		row = []gotgbot.InlineKeyboardButton{
			{Text: NextMonth, CallbackData: fmt.Sprintf("pgm-%d", curMah.AddDate(0, 1, 0).Unix())},
		}
	} else {
		row = []gotgbot.InlineKeyboardButton{
			{Text: NextMonth, CallbackData: fmt.Sprintf("pgm-%d", curMah.AddDate(0, 1, 0).Unix())},
			{Text: PreviousMonth, CallbackData: fmt.Sprintf("pgm-%d", curMah.AddDate(0, -1, 0).Unix())},
		}
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	return markup, nil
}

func createTrainListMarkup(tr database.TrainWR) (*gotgbot.InlineKeyboardMarkup, error) {
	markup := &gotgbot.InlineKeyboardMarkup{}

	trainDayInfo := raja.TrainInfo{
		Source:      raja.Station{Id: tr.Src},
		Destination: raja.Station{Id: tr.Dst},
		ShamsiDate:  ptime.Unix(tr.Day, 0),
	}
	password, err := raja.GetPassword()
	if err != nil {
		return markup, err
	}
	q, err := raja.Encrypt(trainDayInfo.Encode(), password)
	if err != nil {
		return markup, err
	}
	ak, err := raja.GetApiKey()
	if err != nil {
		return markup, err
	}
	trainList, err := raja.GetTrainList(q, &raja.GetTrainListOpt{
		HttpClient: &http.Client{
			Timeout: time.Duration(config.Cfg.Raja.Timeout) * time.Second,
		},
		ApiKey: ak,
	})
	if err != nil {
		if errors.Is(err, raja.ErrTrainsNotFound) {
			markup.InlineKeyboard = append(markup.InlineKeyboard, []gotgbot.InlineKeyboardButton{
				{
					Text:         TrainNotFound,
					CallbackData: "nil",
				},
			})
		}
		if errors.Is(err, raja.ErrGetTrains) || errors.Is(err, raja.ErrBadStatus) || errors.Is(err, raja.ErrGetTrainsDecode) {
			markup.InlineKeyboard = append(markup.InlineKeyboard, []gotgbot.InlineKeyboardButton{
				{
					Text:         RajaErr,
					CallbackData: "nil",
				},
			})
		}
		return markup, err
	}

	for _, train := range trainList.Trains {
		callbackData := fmt.Sprintf("tr-%d-%s", train.RowID, train.ExitTime)
		exitTime, _ := time.ParseInLocation("2006-01-02T15:04:05", train.ExitDateTime, ptime.Iran())
		if time.Now().Unix() >= exitTime.Unix() {
			callbackData = fmt.Sprintf("oldtr-%d", train.RowID)
		}
		markup.InlineKeyboard = append(markup.InlineKeyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf(TrainButtonText, train.ExitTime, strings.TrimSpace(train.WagonName), numToMoney(int(train.Cost/10))),
				CallbackData: callbackData,
			},
		})
	}
	return markup, nil
}

func createListMsg(trs *[]database.TrainWR) string {
	if len(*trs) == 0 {
		return EmptyTrainWR
	}
	msg := ListReqs + "\n"
	for i, tr := range *trs {
		src, _ := Stations.GetPersianName(tr.Src)
		dst, _ := Stations.GetPersianName(tr.Dst)
		msg += fmt.Sprintf(
			"%v.\n"+
				">روز: %v\n"+
				">مبدا: %v\n"+
				">مقصد: %v\n"+
				">ساعت: %v\n"+
				"\n",
			i,
			ptime.Unix(tr.Day, 0).Format(TimeFormat),
			src,
			dst,
			tr.Hour,
		)
	}
	return escapeMarkdown(msg)
}
