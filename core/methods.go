package core

import (
	"RajaBot/config"
	"RajaBot/database"
	"RajaBot/prometheus"
	siteapi "RajaBot/siteApi"
	"RajaBot/siteApi/mrbilit"
	"RajaBot/tools"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FDeghy/RajaGo/raja"
	"github.com/PaulSonOfLars/gotgbot/v2"
	ptime "github.com/yaa110/go-persian-calendar"
)

func sendAlert(trainWR database.TrainWR, train raja.GoTrains) {
	uCacheKey := userCache{
		tgUserId: trainWR.UserID,
		trainId:  trainWR.TrainId,
	}
	mutex.RLock()
	lastUnix, ok := userTimeCache[uCacheKey]
	mutex.RUnlock()
	if !ok {
		lastUnix = 0
	}
	nowUnix := time.Now().Unix()
	if nowUnix-lastUnix < config.Cfg.Raja.AlertEvery {
		return
	}
	src, _ := stations.GetPersianName(trainWR.Src)
	dst, _ := stations.GetPersianName(trainWR.Dst)
	Bot.SendMessage(
		trainWR.UserID,
		fmt.Sprintf(
			AlertMsg,
			train.ExitTime,
			ptime.Unix(trainWR.Day, 0).Format(TrainDate),
			src,
			dst,
			int(train.Counting),
		),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "❌ غیر فعال کردن",
							CallbackData: fmt.Sprintf("canc-%v", trainWR.Id),
						},
					},
					{
						{
							Text: RajaSearchButTxt,
							Url: fmt.Sprintf(
								RajaSearchURL,
								trainWR.Src,
								trainWR.Dst,
								ptime.Unix(trainWR.Day, 0).Format(RajaSearchDateFmt),
							),
						},
					},
				},
			},
		},
	)

	mutex.Lock()
	userTimeCache[uCacheKey] = nowUnix
	prometheus.SetUserTimeCacheCount(len(userTimeCache))
	mutex.Unlock()
}

func sendRtAlert(trainWR database.TrainWR, train siteapi.Train) {
	uCacheKey := userCache{
		tgUserId: trainWR.UserID,
		trainId:  trainWR.TrainId,
	}
	mutex.RLock()
	lastUnix, ok := userTimeCache[uCacheKey]
	mutex.RUnlock()
	if !ok {
		lastUnix = 0
	}
	nowUnix := time.Now().Unix()
	if nowUnix-lastUnix < config.Cfg.Raja.AlertEvery {
		return
	}

	Bot.SendMessage(
		trainWR.UserID,
		fmt.Sprintf(
			AlertMsg,
			train.StartTime,
			ptime.Unix(trainWR.Day, 0).Format(TrainDate),
			train.SourceStationName,
			train.TargetStationName,
			train.SeatRest,
		),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "❌ غیر فعال کردن",
							CallbackData: fmt.Sprintf("canc-%v", trainWR.Id),
						},
					},
				},
			},
		},
	)

	mutex.Lock()
	userTimeCache[uCacheKey] = nowUnix
	prometheus.SetUserTimeCacheCount(len(userTimeCache))
	mutex.Unlock()
}

func sendAlertThrdApp(trainWR database.TrainWR, train *mrbilit.Trains) {
	uCacheKey := userCache{
		tgUserId: trainWR.UserID,
		trainId:  trainWR.TrainId,
	}
	mutex.RLock()
	lastUnix, ok := userTimeCache[uCacheKey]
	mutex.RUnlock()
	if !ok {
		lastUnix = 0
	}
	nowUnix := time.Now().Unix()
	if nowUnix-lastUnix < config.Cfg.Raja.AlertEvery {
		return
	}
	src, _ := stations.GetPersianName(trainWR.Src)
	dst, _ := stations.GetPersianName(trainWR.Dst)
	ti, _ := time.ParseInLocation("2006-01-02T15:04:05", train.DepartureTime, ptime.Iran())
	pt := ptime.New(ti)

	Bot.SendMessage(
		trainWR.UserID,
		fmt.Sprintf(
			AlertMsg,
			pt.Format("HH:mm"),
			pt.Format(TrainDate),
			src,
			dst,
			train.Prices[0].Classes[0].Capacity,
		),
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text:         "❌ غیر فعال کردن",
							CallbackData: fmt.Sprintf("canc-%v", trainWR.Id),
						},
					},
				},
			},
		},
	)

	mutex.Lock()
	userTimeCache[uCacheKey] = nowUnix
	prometheus.SetUserTimeCacheCount(len(userTimeCache))
	mutex.Unlock()
}

func expireWork(oldTrains []*database.TrainWR) {
	for _, i := range oldTrains {
		var src, dst string
		if i.Dst != -1 { // raja api (1)
			src, _ = stations.GetPersianName(i.Src)
			dst, _ = stations.GetPersianName(i.Dst)
		} else { // ticket.rai api (2)
			route := tools.Routes.FindRoute(strconv.Itoa(i.Src))
			s := strings.Split(route.Name, " به ")
			src, dst = s[0], s[1]
		}
		Bot.SendMessage(
			i.UserID,
			fmt.Sprintf(
				ExpireMsg,
				i.Hour,
				ptime.Unix(i.Day, 0).Format(TrainDate),
				src,
				dst,
			),
			nil,
		)
		i.IsDone = true
		database.UpdateTrainWR(i)

		mutex.Lock()
		delete(userTimeCache, userCache{tgUserId: i.UserID, trainId: i.TrainId})
		prometheus.SetUserTimeCacheCount(len(userTimeCache))
		mutex.Unlock()
	}
}

func CancelWork(twrid uint64) error {
	train := database.GetTrainWRByTid(twrid)
	if train == nil {
		return ErrTrainNotFound
	}
	if train.IsDone {
		return ErrTrainAlreadyDone
	}

	var src, dst string
	if train.Dst != -1 { // raja api (1)
		src, _ = stations.GetPersianName(train.Src)
		dst, _ = stations.GetPersianName(train.Dst)
	} else { // ticket.rai api (2)
		route := tools.Routes.FindRoute(strconv.Itoa(train.Src))
		s := strings.Split(route.Name, " به ")
		src, dst = s[0], s[1]
	}
	Bot.SendMessage(
		train.UserID,
		fmt.Sprintf(
			CancelMsg,
			train.Hour,
			ptime.Unix(train.Day, 0).Format(TrainDate),
			src,
			dst,
		),
		nil,
	)
	train.IsDone = true
	database.UpdateTrainWR(train)

	mutex.Lock()
	delete(userTimeCache, userCache{tgUserId: train.UserID, trainId: train.TrainId})
	prometheus.SetUserTimeCacheCount(len(userTimeCache))
	mutex.Unlock()

	return nil
}
