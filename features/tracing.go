package features

import (
	"github.com/keymetrics/pm2-io-apm-go/structures"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"
)

func InitTracing(config *structures.Config) error {
	reporter := zipkinhttpreporter.NewReporter("https://" + config.PublicKey + ":" + config.PublicKey + "@zipkin.cloud.pm2.io/api/v2/spans")

	endpoint, err := openzipkin.NewEndpoint(config.Name, "")
	if err != nil {
		return err
	}

	// OpenCensus
	ze := zipkin.NewExporter(reporter, endpoint)
	trace.RegisterExporter(ze)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	return nil
}
