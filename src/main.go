package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/f-hj/pm2-io-apm-go/services"
	"github.com/pkg/errors"

	"github.com/f-hj/pm2-io-apm-go/structures"

	"github.com/f-hj/pm2-io-apm-go"
)

func main() {
	test := float64(0)
	Pm2Io := pm2_io_apm_go.Pm2Io{}
	Pm2Io.Start("9nc25845w31vqeq", "1e34mwmtaid0pr7", "Golang_application")

	metric := structures.Metric{
		Name:  "test",
		Value: 0,
	}
	services.AddMetric(&metric)

	services.AddAction(&structures.Action{
		ActionName: "Test",
		Callback: func() string {
			log.Println("Action TEST")
			return "Test"
		},
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 1000; i++ {
			fmt.Fprintf(w, "Hello")
		}
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
			Pm2Io.NotifyError(err)
		}
	}()

	go func() {
		ticker := time.NewTicker(4 * time.Second)
		log.Println("created log ticker")
		for {
			<-ticker.C
			//Pm2Io.SendLog("I love logging things")
		}
	}()

	http.ListenAndServe(":8080", nil)
}
