package services_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/services"
)

func TestMetrics(t *testing.T) {
	var handlerCalled = false

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
}
