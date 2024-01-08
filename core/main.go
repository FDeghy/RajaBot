package core

import (
	"RajaBot/config"
	"sync"
)

func StartCore() error {
	workers = make(map[Work]chan struct{})
	mutex = sync.RWMutex{}
	res = make(chan *TrainData, config.Cfg.Raja.Buffer)
	quit := make(chan struct{})
	for i := 0; i < config.Cfg.Raja.Worker; i++ {
		go procWorker(quit)
	}
	return nil
}
