package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/keymetrics/pm2-io-apm-go/structures"
	"github.com/gorilla/websocket"
)

var tempName = ""

type Transporter struct {
	ws *websocket.Conn
	mu sync.Mutex
}

type Message struct {
	Payload interface{} `json:"payload"`
	Channel string      `json:"channel"`
}

func (transporter *Transporter) Connect(publicKey string, privateKey string, name string, version string) {
	u := url.URL{Scheme: "wss", Host: "omicron.keymetrics.io", Path: "/interaction/public"}

	tempName = name

	headers := http.Header{}
	headers.Add("X-KM-PUBLIC", publicKey)
	headers.Add("X-KM-SECRET", privateKey)
	headers.Add("X-KM-SERVER", name)
	headers.Add("X-PM2-VERSION", version)
	headers.Add("X-PROTOCOL-VERSION", "1")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	log.Println("dial")
	if err != nil {
		log.Fatal("dial:", err)
	}

	transporter.ws = c
	go transporter.MessagesHandler()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		log.Println("created ping ticker")
		for {
			<-ticker.C
			transporter.mu.Lock()
			transporter.ws.WriteMessage(websocket.PingMessage, []byte{})
			transporter.mu.Unlock()
		}
	}()
}

func (transporter *Transporter) MessagesHandler() {
	for {
		_, message, err := transporter.ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

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
			name := payload["action_name"]
			switch name {
			case "startLogging":
				transporter.Send("trigger:pm2:result", map[string]interface{}{
					"success": true,
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

	log.Println(string(b))

	transporter.ws.WriteMessage(websocket.TextMessage, b)
}

func (transporter *Transporter) Send(channel string, data interface{}) {
	transporter.SendJson(Message{
		Channel: channel,
		Payload: PayLoad{
			At: time.Now().UnixNano() / int64(time.Millisecond),
			Process: structures.Process{
				PmID:   0,
				Name:   tempName,
				Server: tempName,
			},
			Data:       data,
			Active:     true,
			ServerName: tempName,
			Protected:  false,
			RevCon:     true,
		},
	})
}
