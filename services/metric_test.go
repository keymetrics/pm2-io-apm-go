package services_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/services"
)

func TestMetrics(t *testing.T) {
	t.Run("Generate map", func(t *testing.T) {
		metrics := services.GetMetricsAsMap()
		if metrics == nil {
			t.Fatal("cannot get map")
		}
	})
}
