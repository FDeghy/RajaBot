package rajaHandler

import (
	"RajaBot/config"
	"RajaBot/core"
	"RajaBot/database"
	"RajaBot/siteApi/mrbilit"
	"RajaBot/tools"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
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

// rt - new raja api ticket.rai (2)
func createRoutesMarkup() *gotgbot.InlineKeyboardMarkup {
	markup := &gotgbot.InlineKeyboardMarkup{}
	row := []gotgbot.InlineKeyboardButton{}
	for i, rt := range Routes {
		if i%3 == 0 && i != 0 {
			markup.InlineKeyboard = append(markup.InlineKeyboard, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
		row = append(row, gotgbot.InlineKeyboardButton{
			Text:         rt.Name,
			CallbackData: fmt.Sprintf("rt-%s", rt.ID),
		})
		i++
	}
	if len(row) != 0 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, row)
	}
	return markup
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
	tools.AppendEmptyButton(&row, pad)
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
	tools.AppendEmptyButton(&row, pad)
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
		callbackData := fmt.Sprintf("tr-raja-%d-%s", train.RowID, train.ExitTime)
		exitTime, _ := time.ParseInLocation("2006-01-02T15:04:05", train.ExitDateTime, ptime.Iran())
		if time.Now().Unix() >= exitTime.Unix() {
			callbackData = fmt.Sprintf("oldtr-%d", train.RowID)
		}
		markup.InlineKeyboard = append(markup.InlineKeyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf(TrainButtonText, train.ExitTime, strings.TrimSpace(train.WagonName), tools.NumToMoney(int(train.Cost/10))),
				CallbackData: callbackData,
			},
		})
	}
	return markup, nil
}

func createTrainRtListMarkup(tr database.TrainWR) (*gotgbot.InlineKeyboardMarkup, error) {
	markup := &gotgbot.InlineKeyboardMarkup{}

	route := Routes.FindRoute(strconv.Itoa(tr.Src))
	pt := ptime.Unix(tr.Day, 0)
	pt.At(0, 0, 0, 0)
	// trainList, err := siteapi.GetTrains(route.Src, route.Dst, pt.Format("yyyy/MM/dd"))
	// if err != nil {
	// 	return markup, err
	// }
	trainList, _ := core.UpdateRtsTrains(route, pt)

	maxTries := 1
	for i := 0; i <= maxTries; i++ {
		if i == maxTries && len(trainList) == 0 { // failed (end)
			return &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{
					Text:         TrainNotFound,
					CallbackData: "nil",
				}}},
			}, nil
		} else if len(trainList) == 0 { // try again and try again
			trainList, _ = core.UpdateRtsTrains(route, pt)
		} else { // db have train list
			break
		}
	}

	for _, train := range trainList {
		clock := strings.Split(train.StartTime, ":")
		hour, _ := strconv.Atoi(clock[0])
		minute, _ := strconv.Atoi(clock[1])
		pt.SetHour(hour)
		pt.SetMinute(minute)
		callbackData := fmt.Sprintf("rttr-%d-%s", train.ID, train.StartTime)
		if time.Now().Unix() >= pt.Unix() {
			callbackData = fmt.Sprintf("oldtr-%d", train.ID)
		}
		markup.InlineKeyboard = append(markup.InlineKeyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf(TrainButtonText, train.StartTime, strings.TrimSpace(train.CompartmentHelp+" "+train.CompanyName), tools.NumToMoney(int(train.SeatPrice/10))),
				CallbackData: callbackData,
			},
		})
	}
	return markup, nil
}

func createTrainListThrdAppMarkup(tr database.TrainWR) (*gotgbot.InlineKeyboardMarkup, error) {
	markup := &gotgbot.InlineKeyboardMarkup{}

	date := ptime.Unix(tr.Day, 0).Time().Format("2006-01-02")

	trains, err := mrbilit.GetTrains(strconv.Itoa(tr.Src), strconv.Itoa(tr.Dst), date)
	if err != nil {
		markup.InlineKeyboard = append(markup.InlineKeyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         TrainNotFound,
				CallbackData: "nil",
			},
		})
		return markup, err
	}

	for _, train := range trains {
		ti, _ := time.ParseInLocation("2006-01-02T15:04:05", train.DepartureTime, ptime.Iran())
		exitTime := ptime.New(ti)
		callbackData := fmt.Sprintf("tr-ta-%v-%v", train.ID, exitTime.Format("HH:mm"))
		price := train.Prices[0].Classes[0]
		if time.Now().Unix() >= exitTime.Unix() {
			callbackData = fmt.Sprintf("oldtr-%d", train.ID)
		}
		markup.InlineKeyboard = append(markup.InlineKeyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf(TrainButtonText, exitTime.Format("HH:mm"), fmt.Sprintf("%v %v", train.CorporationName, train.TrainNumber), tools.NumToMoney(int(price.Price/10))),
				CallbackData: callbackData,
			},
		})
	}
	return markup, nil
}

func createListMsg(trs []*database.TrainWR) string {
	if len(trs) == 0 {
		return EmptyTrainWR
	}
	msg := ListReqs + "\n"
	for i, tr := range trs {
		var src, dst string
		if tr.Dst != -1 { // raja api (1)
			src, _ = Stations.GetPersianName(tr.Src)
			dst, _ = Stations.GetPersianName(tr.Dst)
		} else { // ticket.rai api (2)
			route := tools.Routes.FindRoute(strconv.Itoa(tr.Src))
			s := strings.Split(route.Name, " به ")
			src, dst = s[0], s[1]
		}
		msg += fmt.Sprintf(
			"%v\\.\n"+
				">روز\\: %v\n"+
				">مبدا\\: %v\n"+
				">مقصد\\: %v\n"+
				">ساعت\\: %v\n"+
				"\n",
			i+1,
			ptime.Unix(tr.Day, 0).Format(TimeFormat),
			src,
			dst,
			tr.Hour,
		)
	}
	return msg
}

func createListMarkup(trs []*database.TrainWR) *gotgbot.InlineKeyboardMarkup {
	markup := &gotgbot.InlineKeyboardMarkup{}
	if len(trs) == 0 {
		return markup
	}
	for i, tr := range trs {
		markup.InlineKeyboard = append(markup.InlineKeyboard, []gotgbot.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf(CancButtonTxt, i+1),
				CallbackData: fmt.Sprintf("canc-%v", tr.Id),
			},
		})
	}
	return markup
}
