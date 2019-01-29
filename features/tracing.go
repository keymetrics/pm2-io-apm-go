package features

import (
	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/keymetrics/pm2-io-apm-go/structures"
	openzipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"
)

// InitTracing with config and transporter provided
func InitTracing(config *structures.Config, transporter *services.Transporter) error {
	endpoint, err := openzipkin.NewEndpoint(config.Name, "")
	if err != nil {
		return err
	}

	// OpenCensus
	ze := zipkin.NewExporter(NewWsReporter(transporter), endpoint)
	trace.RegisterExporter(ze)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(0.5)})

	return nil
}

// WsReporter is a Zipkin compatible reporter through PM2 WebSocket
type WsReporter struct {
	Transporter *services.Transporter
}

// Send a message using PM2 transporter
func (r *WsReporter) Send(s model.SpanModel) {
	type Alias model.SpanModel

	timers, err := getTimers(s)
	if err != nil {
		return
	}

	t := &struct {
		Process structures.Process `json:"process"`
		Alias
		Timers
	}{
		Process: structures.Process{
			PmID:   0,
			Name:   r.Transporter.Config.Name,
			Server: r.Transporter.Config.ServerName,
		},
		Alias:  (Alias)(s),
		Timers: *timers,
	}
	msg := services.Message{
		Channel: "trace-span",
		Payload: t,
	}
	r.Transporter.SendJson(msg)
}

// Close the reporter (not used, ws handled in transporter)
func (r *WsReporter) Close() error {
	return nil
}

// NewWsReporter create a reporter using specified transporter
func NewWsReporter(transporter *services.Transporter) reporter.Reporter {
	rep := WsReporter{
		Transporter: transporter,
	}
	return &rep
}
