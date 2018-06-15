package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/keymetrics/pm2-io-apm-go/features/metrics"
	"github.com/keymetrics/pm2-io-apm-go/structures"
)

type Transporter struct {
	Config     *structures.Config
	Version    string
	Hostname   string
	ServerName string

	ws         *websocket.Conn
	mu         sync.Mutex
	isHandling bool
}

type Message struct {
	Payload interface{} `json:"payload"`
	Channel string      `json:"channel"`
}

func NewTransporter(config *structures.Config, version string, hostname string, serverName string) *Transporter {
	return &Transporter{
		Config:     config,
		Version:    version,
		Hostname:   hostname,
		ServerName: serverName,

		isHandling: false,
	}
}

func (transporter *Transporter) Connect() {
	u := url.URL{Scheme: "wss", Host: transporter.Config.Server, Path: "/interaction/public"}

	headers := http.Header{}
	headers.Add("X-KM-PUBLIC", transporter.Config.PublicKey)
	headers.Add("X-KM-SECRET", transporter.Config.PrivateKey)
	headers.Add("X-KM-SERVER", transporter.ServerName)
	headers.Add("X-PM2-VERSION", transporter.Version)
	headers.Add("X-PROTOCOL-VERSION", "1")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Println("dial:", err)

		transporter.CloseAndReconnect()
		return
	}

	log.Println("connected")

	transporter.ws = c

	if !transporter.isHandling {
		transporter.SetHandlers()
	}
}

func (transporter *Transporter) SetHandlers() {
	transporter.isHandling = true

	go transporter.MessagesHandler()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			<-ticker.C
			transporter.mu.Lock()
			errPinger := transporter.ws.WriteMessage(websocket.PingMessage, []byte{})
			transporter.mu.Unlock()
			if errPinger != nil {
				log.Println("disconnected pinger:", errPinger)
				transporter.CloseAndReconnect()
				return
			}
		}
	}()
}

func (transporter *Transporter) MessagesHandler() {
	for {
		_, message, err := transporter.ws.ReadMessage()
		if err != nil {
			log.Println("disconnected msgHandler:", err)
			transporter.CloseAndReconnect()
			return
		}

		log.Println(string(message))

		var dat map[string]interface{}

		if err := json.Unmarshal(message, &dat); err != nil {
			panic(err)
		}

		if dat["channel"] == "trigger:action" {
			payload := dat["payload"].(map[string]interface{})
			name := payload["action_name"]

			response := CallAction(name.(string))

			transporter.Send("trigger:action:success", map[string]interface{}{
				"success":     true,
				"id":          payload["process_id"],
				"action_name": name,
			})
			transporter.Send("axm:reply", map[string]interface{}{
				"action_name": name,
				"return":      response,
			})

		} else if dat["channel"] == "trigger:pm2:action" {
			payload := dat["payload"].(map[string]interface{})
			name := payload["method_name"]
			switch name {
			case "startLogging":
				transporter.SendJson(map[string]interface{}{
					"channel": "trigger:pm2:result",
					"payload": map[string]interface{}{
						"ret": map[string]interface{}{
							"err": nil,
						},
					},
				})
				break
			}
		} else {
			log.Println("msg not registered: " + dat["channel"].(string))
		}
	}
}

func (transporter *Transporter) SendJson(msg interface{}) {
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
	}

	transporter.mu.Lock()
	defer transporter.mu.Unlock()

	log.Println("sent")
	transporter.ws.WriteMessage(websocket.TextMessage, b)
}

func (transporter *Transporter) Send(channel string, data interface{}) {
	log.Println("sending:", channel)
	transporter.SendJson(Message{
		Channel: channel,
		Payload: PayLoad{
			At: time.Now().UnixNano() / int64(time.Millisecond),
			Process: structures.Process{
				PmID:   0,
				Name:   transporter.Config.Name,
				Server: transporter.ServerName,
			},
			Data:       data,
			Active:     true,
			ServerName: transporter.ServerName,
			Protected:  false,
			RevCon:     true,
			InternalIP: metrics.LocalIP(),
		},
	})
}

func (transporter *Transporter) CloseAndReconnect() {
	log.Println("CloseAndReconnect")

	log.Println("closing...")
	transporter.ws.Close()
	transporter.Connect()
}
