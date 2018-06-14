package metrics

import (
	"runtime"

	"github.com/keymetrics/pm2-io-apm-go/structures"
)

func GoRoutines() *structures.Metric {
	metric := structures.CreateFuncMetric("GoRoutines", "metric", "routines", func() float64 {
		return float64(runtime.NumGoroutine())
	})
	return &metric
}

func CgoCalls() *structures.Metric {
	last := runtime.NumCgoCall()
	metric := structures.CreateFuncMetric("CgoCalls/sec", "metric", "calls/sec", func() float64 {
		calls := runtime.NumCgoCall()
		v := calls - last
		last = calls
		return float64(v)
	})
	return &metric
}
