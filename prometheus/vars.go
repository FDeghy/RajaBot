package prometheus

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	mutex   = &sync.Mutex{}
	errChan = make(chan error)
)

var (
	fetchWorkersCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "core_fetch_workers_count",
		Help: "len of core.workers",
	})
	userTimeCacheCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "core_user_time_cache_count",
		Help: "len of core.userTimeCache",
	})
)
