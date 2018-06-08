package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/f-hj/pm2-io-apm-go"
)

func main() {
	test := float64(0)
	Pm2Io := pm2_io_apm_go.Pm2Io{
		PublicKey:  "9nc25845w31vqeq",
		PrivateKey: "1e34mwmtaid0pr7",
		Name:       "Golang_application",
		AxmActions: []pm2_io_apm_go.AxmAction{
			pm2_io_apm_go.AxmAction{
				ActionName: "Test",
				Callback: func() string {
					log.Println("THIS IS TEST ACTION")
					return "YES"
				},
			},
		},
		AxmMonitor: map[string]pm2_io_apm_go.AxmMonitor{
			"Test": pm2_io_apm_go.AxmMonitor{
				Value: test,
			},
		},
	}

	Pm2Io.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 1000; i++ {
			fmt.Fprintf(w, "Hello")
		}
	})

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		log.Println("created ticker")
		for {
			<-ticker.C
			test++
			Pm2Io.SetProbe("Test", test)
			Pm2Io.NotifyError(errors.New("Niaha"))
		}
	}()

	http.ListenAndServe(":8080", nil)
}
