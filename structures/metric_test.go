package structures_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/structures"
)

func MetricTest(t *testing.T) {
	var funcMetric structures.Metric
	var stdMetric structures.Metric

	t.Run("create func metric", func(t *testing.T) {
		funcMetric = structures.CreateFuncMetric("name", "category", "unit", func() float64 {
			return 42
		})
	})

	t.Run("get func value", func(t *testing.T) {
		val := funcMetric.Get()
		if val != 42 {
			t.Fatal("bad func ret value")
		}
	})

	t.Run("create std metric", func(t *testing.T) {
		stdMetric = structures.CreateMetric("name2", "category2", "unit2")
	})

	t.Run("set value", func(t *testing.T) {
		stdMetric.Set(20000)
	})

	t.Run("get value", func(t *testing.T) {
		val := stdMetric.Get()
		if val != 20000 {
			t.Fatal("bad std value")
		}
	})
}
