package database

import (
	siteapi "RajaBot/siteApi"
	"encoding/json"
	"time"

	ptime "github.com/yaa110/go-persian-calendar"
)

func NewSubscription(userId int64) *Subscription {
	return &Subscription{
		UserID:         userId,
		IsTrial:        false,
		ExpirationDate: 0,
		IsEnabled:      false,
		RegisteryDate:  time.Now().Unix(),
	}
}

func NewTgUser(userId int64) *TgUser {
	return &TgUser{
		UserID: userId,
		IsVip:  false,
		State:  "normal",
	}
}

func NewRTTrain(src, dst string, date ptime.Time, trains []siteapi.Train) *RTTrain {
	pt := date
	pt.At(0, 0, 0, 0)
	jsonTrains, _ := json.Marshal(trains)
	rts := &RTTrain{
		Src:    src,
		Dst:    dst,
		Date:   pt.Unix(),
		Trains: string(jsonTrains),
	}

	return rts
}

func (r *RTTrain) GetTrains() []siteapi.Train {
	var trains []siteapi.Train
	err := json.Unmarshal([]byte(r.Trains), &trains)
	if err != nil {
		return nil
	}

	return trains
}

func (r *RTTrain) SetTrains(trs []siteapi.Train) {
	jsonTrains, _ := json.Marshal(trs)
	r.Trains = string(jsonTrains)
}
