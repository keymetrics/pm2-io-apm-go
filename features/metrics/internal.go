package metrics

import (
	"runtime"

	"github.com/keymetrics/pm2-io-apm-go/structures"
)

type MetricsMemStats struct {
	Initied   bool
	NumGC     *structures.Metric
	LastNumGC float64

	NumMallocs     *structures.Metric
	LastNumMallocs float64

	NumFree     *structures.Metric
	LastNumFree float64

	HeapAlloc *structures.Metric

	Pause     *structures.Metric
	LastPause float64
}

var GlobalMetricsMemStats MetricsMemStats

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

func InitMetricsMemStats() {
	numGC := structures.CreateMetric("GCRuns/sec", "metric", "runs")
	numMalloc := structures.CreateMetric("mallocs/sec", "metric", "mallocs")
	numFree := structures.CreateMetric("free/sec", "metric", "frees")
	heapAlloc := structures.CreateMetric("heapAlloc", "metric", "bytes")
	pause := structures.CreateMetric("Pause/sec", "metric", "ns/sec")

	GlobalMetricsMemStats = MetricsMemStats{
		Initied:    true,
		NumGC:      &numGC,
		NumMallocs: &numMalloc,
		NumFree:    &numFree,
		HeapAlloc:  &heapAlloc,
		Pause:      &pause,
	}
}

func Handler() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	GlobalMetricsMemStats.NumGC.Set(float64(stats.NumGC) - GlobalMetricsMemStats.LastNumGC)
	GlobalMetricsMemStats.LastNumGC = float64(stats.NumGC)

	GlobalMetricsMemStats.NumMallocs.Set(float64(stats.Mallocs) - GlobalMetricsMemStats.LastNumMallocs)
	GlobalMetricsMemStats.LastNumMallocs = float64(stats.Mallocs)

	GlobalMetricsMemStats.NumFree.Set(float64(stats.Frees) - GlobalMetricsMemStats.LastNumFree)
	GlobalMetricsMemStats.LastNumFree = float64(stats.Frees)

	GlobalMetricsMemStats.HeapAlloc.Set(float64(stats.HeapAlloc))

	GlobalMetricsMemStats.Pause.Set(float64(stats.PauseTotalNs) - GlobalMetricsMemStats.LastPause)
	GlobalMetricsMemStats.LastPause = float64(stats.PauseTotalNs)
}
