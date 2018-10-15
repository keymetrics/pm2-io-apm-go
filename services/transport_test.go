package services_test

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/keymetrics/pm2-io-apm-go/features/metrics"
	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/keymetrics/pm2-io-apm-go/structures"
	gock "gopkg.in/h2non/gock.v1"
)

var nbConnected int

func TestTransport(t *testing.T) {
	var wssServer *httptest.Server

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

	t.Run("Mock wss", func(t *testing.T) {
		wssServer = httptest.NewServer(http.HandlerFunc(echo))
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
				Hostname:    transporter.ServerName,
			},
		}

		gock.New("https://root.keymetrics.io").
			Post("/api/node/verifyPM2").
			MatchType("json").
			JSON(verify).
			Reply(200).
			JSON(services.VerifyResponse{
				Endpoints: services.Endpoints{
					WS: "ws" + strings.TrimPrefix(wssServer.URL, "http"),
				},
			})
	})

	t.Run("Connect transporter", func(t *testing.T) {
		transporter.Connect()
		if !transporter.IsConnected() {
			t.Fatal("transporter is not connected")
		}
	})

	t.Run("Wait for ws connection", func(t *testing.T) {
		if nbConnected == 0 {
			t.Fatal("WS not connected")
		}
		if nbConnected != 1 {
			t.Fatal("WS connected more than one time")
		}
	})
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	nbConnected++
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
	nbConnected--
}
