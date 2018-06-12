package features

import (
	"fmt"
	"log"

	"github.com/f-hj/pm2-io-apm-go/services"
	"github.com/pkg/errors"
)

type Notifier struct {
	Transporter *services.Transporter
}

type Error struct {
	Message string `json:"message"`
	Stack   string `json:"stack"`
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (notifier *Notifier) Error(err error) {
	log.Println("ERRRRROOOOOOORRRR")

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

func (notifier *Notifier) Log(log string) {
	notifier.Transporter.Send("log", log)
}
