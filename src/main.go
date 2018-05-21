package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/f-hj/pm2-io-apm-go"
)

func main() {
	Pm2Io := pm2_io_apm_go.Pm2Io{
		PublicKey:  "9nc25845w31vqeq",
		PrivateKey: "1e34mwmtaid0pr7",
		Name:       "Golang_connector",
		AxmActions: []pm2_io_apm_go.AxmAction{
			pm2_io_apm_go.AxmAction{
				ActionName: "Test",
				Callback: func() string {
					log.Println("THIS IS TEST ACTION")
					return "YES"
				},
			},
		},
	}

	Pm2Io.Start()
	Pm2Io.SendStatus()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 1000; i++ {
			fmt.Fprintf(w, "Hello")
		}
	})

	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		Pm2Io.SendStatus()
	}
}
