package structures

// Metric like AxmMonitor
type Metric struct {
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	Category string  `json:"type"`
	Historic bool    `json:"historic"`
	Unit     string  `json:"unit"`
}

func (metric *Metric) Set(value float64) {
	metric.Value = value
}

func CreateMetric(name string, category string, unit string) Metric {
	return Metric{
		Name:     name,
		Category: category,
		Unit:     unit,
		Value:    0,

		Historic: true,
	}
}
