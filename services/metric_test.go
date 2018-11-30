package services_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/structures"

	"github.com/keymetrics/pm2-io-apm-go/services"
)

func TestMetrics(t *testing.T) {
	var handlerCalled = false
	stdMetric := structures.CreateMetric("name", "category", "unit")

	t.Run("Generate map", func(t *testing.T) {
		metrics := services.GetMetricsAsMap()
		if metrics == nil {
			t.Fatal("cannot get map")
		}
	})

	t.Run("Can add handler", func(t *testing.T) {
		services.AttachHandler(func() {
			handlerCalled = true
		})
	})

	t.Run("Check if handler called", func(t *testing.T) {
		metrics := services.GetMetricsAsMap()
		if metrics == nil {
			t.Fatal("cannot get map")
		}

		if !handlerCalled {
			t.Fatal("Handler attached but never called")
		}
	})

	t.Run("add metric", func(t *testing.T) {
		services.AddMetric(&stdMetric)
	})

	t.Run("should be in array", func(t *testing.T) {
		if len(services.Metrics) != 1 {
			t.Fatal("wanted 1 metric, got " + string(len(services.Metrics)))
		}
	})

	t.Run("should be in map", func(t *testing.T) {
		metrics := services.GetMetricsAsMap()
		if metrics == nil {
			t.Fatal("cannot get map")
		}
		if metrics["name"] == nil {
			t.Fatal("not inserted")
		}
		if metrics["name"].Unit != "unit" {
			t.Fatal("metric data problem")
		}
	})
}
