package core

import (
	"RajaBot/config"
	"RajaBot/database"
	"RajaBot/prometheus"
	siteapi "RajaBot/siteApi"
	"RajaBot/tools"
	"strconv"
	"strings"
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
			prometheus.SetFetchWorkersCount(len(workers))
			mutex.Unlock()
			return
		}
	}
}

func procWorker(q chan struct{}) {
	for {
		select {
		case data := <-res: // raja api (1)
			trWRs := database.GetActiveTrainWRsByInfo(data.Work.Day, data.Work.Src, data.Work.Dst)
			if len(trWRs) == 0 {
				mutex.RLock()
				close(workers[data.Work])
				mutex.RUnlock()
				continue
			}
			for _, tr := range data.TrainList.Trains {
				trExitTime, _ := time.ParseInLocation("2006-01-02T15:04:05", tr.ExitDateTime, ptime.Iran())
				trWR := database.FilterTrainWRsByTrainId(tr.RowID, trWRs)
				if len(trWR) == 0 {
					continue
				}
				if time.Now().Unix() >= trExitTime.Unix() {
					expireWork(trWR)
					continue
				}
				if tr.Counting > 0 {
					for _, trWRData := range trWR {
						// shayad ham go sendAlert
						sendAlert(*trWRData, tr)
					}
				}
			}
		case data := <-rtRes: // ticket.rai api (2)
			trWRs := database.GetActiveTrainWRsByInfo(data.Work.Day, data.Work.Src, data.Work.Dst)
			if len(trWRs) == 0 {
				mutex.RLock()
				close(workers[data.Work])
				mutex.RUnlock()
				continue
			}
			for _, tr := range data.TrainList {
				pt := ptime.Unix(data.Work.Day, 0)
				clock := strings.Split(tr.StartTime, ":")
				hour, _ := strconv.Atoi(clock[0])
				minute, _ := strconv.Atoi(clock[1])
				pt.SetHour(hour)
				pt.SetMinute(minute)
				trWR := database.FilterTrainWRsByTrainId(tr.ID, trWRs)
				if len(trWR) == 0 {
					continue
				}
				if time.Now().Unix() >= pt.Unix() {
					expireWork(trWR)
					continue
				}
				if tr.SeatRest > 0 {
					for _, trWRData := range trWR {
						// shayad ham go sendAlert
						sendRtAlert(*trWRData, *tr)
					}
				}
			}
		case <-q:
			return
		}
	}
}

func rtFetchWorker(wk Work, q chan struct{}) {
	ticker := time.NewTicker(time.Duration(config.Cfg.Raja.CheckEvery) * time.Second)

	route := tools.Routes.FindRoute(strconv.Itoa(wk.Src))
	pt := ptime.Unix(wk.Day, 0)

	for {
		select {
		case <-ticker.C:
			trainList, err := siteapi.GetTrains(route.Src, route.Dst, pt.Format("yyyy/MM/dd"))
			if err == nil {
				rtRes <- &RtTrainData{
					TrainList: trainList,
					Work:      wk,
				}
			}
		case <-q:
			ticker.Stop()
			mutex.Lock()
			delete(workers, wk)
			prometheus.SetFetchWorkersCount(len(workers))
			mutex.Unlock()
			return
		}
	}
}
