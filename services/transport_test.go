package services_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/keymetrics/pm2-io-apm-go/features/metrics"
	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/keymetrics/pm2-io-apm-go/structures"
	gock "gopkg.in/h2non/gock.v1"
)

var nbConnected int

func TestTransport(t *testing.T) {
	var wssServer *httptest.Server
	var gockVerif *gock.Response

	var transporter *services.Transporter
	defaultNode := "api.cloud.pm2.io"

	t.Run("Create transporter", func(t *testing.T) {
		transporter = services.NewTransporter(&structures.Config{
			PublicKey:  "pubKey",
			PrivateKey: "privKey",
			Node:       &defaultNode, // Normally set by pm2io
		}, "version")

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
				MachineName: transporter.Config.Hostname,
				Cpus:        runtime.NumCPU(),
				Memory:      metrics.TotalMem(),
				Pm2Version:  transporter.Version,
				Hostname:    transporter.Config.ServerName,
			},
		}

		gockVerif = gock.New("https://api.cloud.pm2.io").
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
		if nbConnected != 1 {
			t.Fatal("WS connected wanted: 1, connected: " + strconv.Itoa(nbConnected))
		}
	})

	t.Run("Try to close and reconnect", func(t *testing.T) {
		transporter.CloseAndReconnect()
		time.Sleep(2 * time.Second)
		if nbConnected != 1 {
			t.Fatal("WS connected wanted: 1, connected: " + strconv.Itoa(nbConnected))
		}
	})

	t.Run("Shouldn't crash without WSS", func(t *testing.T) {
		wssServer.Close()
		nbConnected = 0
		time.Sleep(2 * time.Second)
	})

	t.Skip("Should get new node and connect to it")
	t.Run("Should get new node and connect to it", func(t *testing.T) {
		wssServer = httptest.NewServer(http.HandlerFunc(echo))

		gockVerif.JSON(services.VerifyResponse{
			Endpoints: services.Endpoints{
				WS: "ws" + strings.TrimPrefix(wssServer.URL, "http"),
			},
		})

		newNode := "ws" + strings.TrimPrefix(wssServer.URL, "http")
		transporter.Config.Node = &newNode
		time.Sleep(8 * time.Second)
		if nbConnected != 1 {
			t.Log("Transporter node: " + *transporter.Config.Node)
			t.Log("Node wanted: " + "ws" + strings.TrimPrefix(wssServer.URL, "http"))
			t.Fatal("WS connected wanted: 1, connected: " + strconv.Itoa(nbConnected))
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
	log.Println("New conn")
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("Conn closed")
			nbConnected--
			return
		}
	}
}
