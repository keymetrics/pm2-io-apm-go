package metrics_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/features/metrics"
)

func TestSystem(t *testing.T) {
	t.Run("Get cpu percent value", func(t *testing.T) {
		value, err := metrics.CPUPercent()
		if err != nil {
			t.Fatal(err)
		}
		if value < 0 {
			t.Fatal("value less than 0")
		}
	})

	t.Run("Get cpu name", func(t *testing.T) {
		value := metrics.CPUName()
		if len(value) == 0 {
			t.Fatal("name empty")
		}
	})

	t.Run("Get cpu load", func(t *testing.T) {
		value := metrics.CPULoad()
		if len(value) != 3 {
			t.Fatal("load empty")
		}
	})

	t.Run("Get local ip", func(t *testing.T) {
		value := metrics.LocalIP()
		if len(value) == 0 {
			t.Fatal("ip empty")
		}
	})
}
