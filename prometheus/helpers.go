package prometheus

import (
	"RajaBot/config"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartProm() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(config.Cfg.Prometheus.AddressBind, mux)
		errChan <- err
	}()
	select {
	case err := <-errChan:
		return err
	case <-time.After(3 * time.Second):
		break
	}
	log.Println("Prometheus -> metrics http server started.")
	return nil
}

func SetFetchWorkersCount(i int) {
	mutex.Lock()
	fetchWorkersCount.Set(float64(i))
	mutex.Unlock()
}

func SetUserTimeCacheCount(i int) {
	mutex.Lock()
	userTimeCacheCount.Set(float64(i))
	mutex.Unlock()
}
