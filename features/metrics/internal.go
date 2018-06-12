package metrics

import (
	"runtime"

	"github.com/keymetrics/pm2-io-apm-go/structures"
)

func GoRoutines() *structures.Metric {
	return &structures.Metric{
		Name: "GoRoutines",
		Function: func() float64 {
			return float64(runtime.NumGoroutine())
		},
	}
}

func CGoCalls() *structures.Metric {
	return &structures.Metric{
		Name: "CGoCalls",
		Function: func() float64 {
			return float64(runtime.NumCgoCall())
		},
	}
}
