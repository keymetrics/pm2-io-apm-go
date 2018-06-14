package structures

// Metric like AxmMonitor
type Metric struct {
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	Category string  `json:"type"`
	Historic bool    `json:"historic"`
	Unit     string  `json:"unit"`

	Aggregation string `json:"agg_type"`

	Function func() float64 `json:"-"`
}

func (metric *Metric) Get() float64 {
	if metric.Function != nil {
		metric.Value = metric.Function()
	}
	return metric.Value
}

func (metric *Metric) Set(value float64) {
	metric.Value = value
}

var defaultAggregation = "avg"

func CreateMetric(name string, category string, unit string) Metric {
	return Metric{
		Name:        name,
		Category:    category,
		Unit:        unit,
		Value:       0,
		Aggregation: defaultAggregation,

		Historic: true,
	}
}

func CreateFuncMetric(name string, category string, unit string, cb func() float64) Metric {
	return Metric{
		Name:        name,
		Category:    category,
		Unit:        unit,
		Aggregation: defaultAggregation,

		Historic: true,
		Function: cb,
	}
}
