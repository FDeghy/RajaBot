package core

import (
	"RajaBot/config"
	"RajaBot/database"
	"log"
	"sync"

	"github.com/FDeghy/RajaGo/raja"
)

func StartCore() error {
	workers = make(map[Work]chan struct{})
	userTimeCache = make(map[int64]int64)
	mutex = sync.RWMutex{}
	res = make(chan *TrainData, config.Cfg.Raja.Buffer)
	sts, err := raja.GetStations()
	if err != nil {
		return nil
	}
	stations = sts

	quit := make(chan struct{})
	for i := 0; i < config.Cfg.Raja.Worker; i++ {
		go procWorker(quit)
	}
	log.Printf("Core -> %v procWorker started.", config.Cfg.Raja.Worker)

	uncompTrainWRs := database.GetAllActiveTrainWRs()
	for _, i := range *uncompTrainWRs {
		HandleGoFetch(&i)
	}
	log.Printf("Core -> %v uncompleted task sent to handler.", len(*uncompTrainWRs))

	return nil
}
