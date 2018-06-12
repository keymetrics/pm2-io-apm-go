package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/pkg/errors"

	"github.com/keymetrics/pm2-io-apm-go/structures"

	"github.com/keymetrics/pm2-io-apm-go"
)

func main() {
	test := float64(0)
	Pm2Io := pm2io.Pm2Io{}
	Pm2Io.Start("9nc25845w31vqeq", "1e34mwmtaid0pr7", "Golang_application")

	metric := structures.CreateMetric("test", "metric", "unit")
	services.AddMetric(&metric)

	nbreq := structures.Metric{
		Name:  "nbreq",
		Value: 0,
	}
	services.AddMetric(&nbreq)

	services.AddAction(&structures.Action{
		ActionName: "Test",
		Callback: func() string {
			log.Println("Action TEST")
			return "I am the test answer"
		},
	})

	services.AddAction(&structures.Action{
		ActionName: "Get env",
		Callback: func() string {
			return strings.Join(os.Environ(), "\n")
		},
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 1000; i++ {
			fmt.Fprintf(w, "Hello")
		}
		nbreq.Value++
	})

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		log.Println("created 2s ticker")
		for {
			<-ticker.C
			test++
			metric.Set(test)
			cause := errors.New("Niaha")
			err := errors.WithStack(cause)
			Pm2Io.Notifier.Error(err)
		}
	}()

	go func() {
		ticker := time.NewTicker(4 * time.Second)
		log.Println("created log ticker")
		for {
			<-ticker.C
			Pm2Io.Notifier.Log("I love logging things")
		}
	}()

	/*go func() {
		ticker := time.NewTicker(6 * time.Second)
		log.Println("created log ticker")
		for {
			<-ticker.C
			cause := errors.New("Fatal panic error")
			err := errors.WithStack(cause)
			Pm2Io.Panic(err)
		}
	}()*/

	http.ListenAndServe(":8080", nil)
}
