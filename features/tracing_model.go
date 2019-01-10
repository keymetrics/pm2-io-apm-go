package features

import (
	"errors"
	"time"

	"github.com/openzipkin/zipkin-go/model"
)

// Timers store timestamp and duration in JSON format from model.SpanModel
type Timers struct {
	T int64 `json:"timestamp,omitempty"`
	D int64 `json:"duration,omitempty"`
}

func getTimers(s model.SpanModel) (*Timers, error) {
	var timestamp int64
	if !s.Timestamp.IsZero() {
		if s.Timestamp.Unix() < 1 {
			// Zipkin does not allow Timestamps before Unix epoch
			return nil, errors.New("ErrValidTimestampRequired")
		}
		timestamp = s.Timestamp.Round(time.Microsecond).UnixNano() / 1e3
	}

	if s.Duration < time.Microsecond {
		if s.Duration < 0 {
			// negative duration is not allowed and signals a timing logic error
			return nil, errors.New("ErrValidDurationRequired")
		} else if s.Duration > 0 {
			// sub microsecond durations are reported as 1 microsecond
			s.Duration = 1 * time.Microsecond
		}
	} else {
		// Duration will be rounded to nearest microsecond representation.
		//
		// NOTE: Duration.Round() is not available in Go 1.8 which we still support.
		// To handle microsecond resolution rounding we'll add 500 nanoseconds to
		// the duration. When truncated to microseconds in the call to marshal, it
		// will be naturally rounded. See TestSpanDurationRounding in span_test.go
		s.Duration += 500 * time.Nanosecond
	}

	if s.LocalEndpoint.Empty() {
		s.LocalEndpoint = nil
	}

	if s.RemoteEndpoint.Empty() {
		s.RemoteEndpoint = nil
	}

	timers := Timers{
		T: timestamp,
		D: s.Duration.Nanoseconds() / 1e3,
	}

	return &timers, nil
}
