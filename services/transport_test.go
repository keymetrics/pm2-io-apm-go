package services_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/keymetrics/pm2-io-apm-go/features/metrics"
	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/keymetrics/pm2-io-apm-go/structures"
	gock "gopkg.in/h2non/gock.v1"
)

func TestTransport(t *testing.T) {
	var transporter *services.Transporter

	t.Run("Create transporter", func(t *testing.T) {
		transporter = services.NewTransporter(&structures.Config{
			PublicKey:  "pubKey",
			PrivateKey: "privKey",
		}, "version", "hostname", "serverName", "root.keymetrics.io")

		if transporter == nil {
			t.Fatal("transporter is nil")
		}
	})

	t.Run("Connect transporter", func(t *testing.T) {
		transporter.Connect()
		if transporter.IsConnected() {
			t.Fatal("transporter shouldn't be connected")
		}
	})

	t.Run("Mock http server for verifier", func(t *testing.T) {
		verify := services.Verify{
			PublicId:  transporter.Config.PublicKey,
			PrivateId: transporter.Config.PrivateKey,
			Data: services.VerifyData{
				MachineName: transporter.ServerName,
				Cpus:        runtime.NumCPU(),
				Memory:      metrics.TotalMem(),
				Pm2Version:  transporter.Version,
				Hostname:    transporter.Hostname,
			},
		}

		gock.New("https://root.keymetrics.io").
			Post("/api/node/verifyPM2").
			MatchType("json").
			JSON(verify).
			Reply(200).
			JSON(services.VerifyResponse{
				Endpoints: services.Endpoints{
					WS: "ws://127.0.0.1/myWs",
				},
			})
	})

	t.Run("Wait for node", func(t *testing.T) {
		t.Log("start loop")
		for i := 0; i < 20; i++ {
			srv := transporter.GetServer()
			t.Log(srv)
			if srv != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
		t.Fatal("Server not detected")
	})

	t.Run("Wait for ws connection", func(t *testing.T) {

	})
}
