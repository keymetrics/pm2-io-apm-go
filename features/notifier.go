package features

import (
	"fmt"

	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/pkg/errors"
)

// Notifier with transporter
type Notifier struct {
	Transporter *services.Transporter
}

// Error packet to KM
type Error struct {
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Error to KM
func (notifier *Notifier) Error(err error) {
	stack := ""
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			stack += fmt.Sprintf("%+v", f)
		}
	} else {
		stack = fmt.Sprintf("%+v", err)
	}
	notifier.Transporter.Send("process:exception", Error{
		Message: err.Error(),
		Stack:   stack,
	})
}

// Log to KM
func (notifier *Notifier) Log(log string) {
	notifier.Transporter.Send("logs", log)
}
