package services

import "github.com/keymetrics/pm2-io-apm-go/structures"

var Metrics []*structures.Metric

func AddMetric(metric *structures.Metric) {
	Metrics = append(Metrics, metric)
}

func GetMetricsAsMap() map[string]*structures.Metric {
	m := make(map[string]*structures.Metric, len(Metrics))
	for _, metric := range Metrics {
		metric.Get()
		m[metric.Name] = metric
	}
	return m
}
