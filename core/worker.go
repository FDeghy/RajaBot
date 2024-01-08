package core

import (
	"RajaBot/config"
	"RajaBot/database"
	"time"

	"github.com/FDeghy/RajaGo/raja"
	ptime "github.com/yaa110/go-persian-calendar"
)

func fetchWorker(wk Work, q chan struct{}, query raja.Query, opt *raja.GetTrainListOpt) {
	ticker := time.NewTicker(time.Duration(config.Cfg.Raja.CheckEvery) * time.Second)

	for {
		select {
		case <-ticker.C:
			trainList, err := raja.GetTrainList(query, opt)
			if err == nil {
				res <- &TrainData{
					TrainList: trainList,
					Work:      wk,
				}
			}
		case <-q:
			ticker.Stop()
			mutex.Lock()
			delete(workers, wk)
			mutex.Unlock()
			return
		}
	}
}

func procWorker(q chan struct{}) {
	var data *TrainData
	for {
		select {
		case data = <-res:
			trWRs := database.GetActiveTrainWRsByInfo(data.Work.Day, data.Work.Src, data.Work.Dst)
			if len(*trWRs) == 0 {
				mutex.RLock()
				close(workers[data.Work])
				mutex.RUnlock()
				continue
			}
			for _, tr := range data.TrainList.Trains {
				trExitTime, _ := time.ParseInLocation("2006-01-02T15:04:05", tr.ExitDateTime, ptime.Iran())
				if time.Now().Unix() >= trExitTime.Unix() {
					expireWork(tr.RowID)
					continue
				}
				if tr.Counting > 0 {
					trWR := database.GetActiveTrainWRsByTrainId(tr.RowID)
					if len(*trWR) == 0 {
						continue
					}
					for _, trWRData := range *trWR {
						// shayad ham go sendAlert
						sendAlert(&trWRData, &tr)
					}
				}
			}
		case <-q:
			return
		}
	}
}
