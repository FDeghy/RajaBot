package core

import (
	"RajaBot/config"
	"RajaBot/database"
	"net/http"
	"time"

	"github.com/FDeghy/RajaGo/raja"
	ptime "github.com/yaa110/go-persian-calendar"
)

func HandleGoFetch(user *database.TgUser, tr *database.TrainWR) error {
	wk := Work{
		Src: tr.Src,
		Dst: tr.Dst,
		Day: tr.Day,
	}
	mutex.RLock()
	_, ok := workers[wk]
	mutex.RUnlock()
	if ok {
		return nil
	}
	// create fetchWorker
	trainDayInfo := raja.TrainInfo{
		Source:      raja.Station{Id: tr.Src},
		Destination: raja.Station{Id: tr.Dst},
		ShamsiDate:  ptime.Unix(tr.Day, 0),
	}
	password, err := raja.GetPassword()
	if err != nil {
		return err
	}
	query, err := raja.Encrypt(trainDayInfo.Encode(), password)
	if err != nil {
		return err
	}
	ak, err := raja.GetApiKey()
	if err != nil {
		return err
	}
	opt := &raja.GetTrainListOpt{
		HttpClient: &http.Client{
			Timeout: time.Duration(config.Cfg.Raja.Timeout) * time.Second,
		},
		ApiKey: ak,
	}
	quit := make(chan struct{})
	go fetchWorker(wk, quit, query, opt)
	mutex.Lock()
	workers[wk] = quit
	mutex.Unlock()
	return nil
}
