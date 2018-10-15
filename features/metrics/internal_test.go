package metrics_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/structures"

	"github.com/keymetrics/pm2-io-apm-go/features/metrics"
)

func TestInternal(t *testing.T) {
	var goRoutines *structures.Metric
	t.Run("Get correct structures of goroutines", func(t *testing.T) {
		goRoutines = metrics.GoRoutines()
	})

	t.Run("Should get a value", func(t *testing.T) {
		v := goRoutines.Get()
		if v == 0 {
			t.Fatal("Metric shouldn't be 0")
		}
	})

	t.Run("Should init metrics stats", func(t *testing.T) {
		metrics.InitMetricsMemStats()
	})

	t.Run("should run handler without problem", func(t *testing.T) {
		metrics.Handler()
	})
}
