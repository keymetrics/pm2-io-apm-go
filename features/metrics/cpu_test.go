package metrics_test

import (
	"testing"

	"github.com/f-hj/pm2-io-apm-go/features/metrics"
)

func TestCpu(t *testing.T) {
	t.Run("Get value", func(t *testing.T) {
		value, err := metrics.CPUPercent()
		if err != nil {
			t.Fatal(err)
		}
		if value < 0 {
			t.Fatal("value less than 0")
		}
	})
}
