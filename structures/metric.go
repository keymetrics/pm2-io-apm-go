package structures

// Metric like AxmMonitor
type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func (metric *Metric) Set(value float64) {
	metric.Value = value
}
