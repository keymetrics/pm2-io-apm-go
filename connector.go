package pm2_io_apm_go

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/process"
)

type Pm2Io struct {
	PublicKey  string
	PrivateKey string
	Name       string
	AxmActions []AxmAction

	ws           *websocket.Conn
	startTime    time.Time
	lastCpuTotal float64
}

func (pm2io *Pm2Io) Start() {
	pm2io.AxmActions = append(pm2io.AxmActions, AxmAction{
		ActionName: "km:heapdump",
		Callback: func() string {
			log.Println("MEMORY PROFIIIIIIIIILING")
			return ""
		},
	})
	pm2io.AxmActions = append(pm2io.AxmActions, AxmAction{
		ActionName: "km:cpuprofiling:start",
		Callback: func() string {
			log.Println("CPUUUUUUUU PROFIIIIIIIIILING start")
			return ""
		},
	})
	pm2io.AxmActions = append(pm2io.AxmActions, AxmAction{
		ActionName: "km:cpuprofiling:stop",
		Callback: func() string {
			log.Println("CPUUUUUUUU PROFIIIIIIIIILING stop")
			return ""
		},
	})

	pm2io.startTime = time.Now()
	u := url.URL{Scheme: "wss", Host: "staging.keymetrics.io", Path: "/interaction/public"}

	headers := http.Header{}
	headers.Add("X-KM-PUBLIC", pm2io.PublicKey)
	headers.Add("X-KM-DATA", "")
	headers.Add("X-KM-SERVER", pm2io.Name)
	headers.Add("X-PM2-VERSION", "10.0.0")
	headers.Add("X-PROTOCOL-VERSION", "1")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	log.Println("dial")
	if err != nil {
		log.Fatal("dial:", err)
	}

	pm2io.ws = c

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var dat map[string]interface{}

			if err := json.Unmarshal(message, &dat); err != nil {
				panic(err)
			}
			fmt.Println(dat)
			if dat["channel"] == "trigger:action" || dat["channel"] == "trigger:scoped_action" {
				payload := dat["payload"].(map[string]interface{})
				name := payload["action_name"]
				for _, i := range pm2io.AxmActions {
					if i.ActionName == name {
						response := i.Callback()
						res := MessageMap{
							Channel: "trigger:action:success",
							Payload: map[string]interface{}{
								"success":     true,
								"id":          payload["process_id"],
								"action_name": name,
								"return":      response,
							},
						}
						jsonString, _ := json.Marshal(res)
						pm2io.ws.WriteMessage(websocket.TextMessage, jsonString)
					}
				}
			}
		}
	}()
}

func (pm2io *Pm2Io) CPUPercent() (float64, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, err
	}
	crt_time, err := p.CreateTime()
	if err != nil {
		return 0, err
	}

	cput, err := p.Times()
	if err != nil {
		return 0, err
	}

	created := time.Unix(0, crt_time*int64(time.Millisecond))
	totalTime := time.Since(created).Seconds()
	log.Println("totalTime     ", totalTime)
	log.Println("cpuTotal      ", cput.Total())
	log.Println("lastCpuTota   ", pm2io.lastCpuTotal)

	if totalTime <= 0 {
		return 0, nil
	}

	val := (cput.Total() - pm2io.lastCpuTotal) * 100
	pm2io.lastCpuTotal = cput.Total()

	log.Println("return        ", val)

	return val, nil
}

func (pm2io *Pm2Io) SendStatus() {
	p, err := process.NewProcess(int32(os.Getpid()))
	cp, _ := pm2io.CPUPercent()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	kmProc := []Process{}

	options := AxmOptions{
		HeapDump:  true,
		Profiling: true,
	}

	parent, _ := p.Parent()
	sProc, _ := parent.Children()
	for _, p := range sProc {
		kmProc = append(kmProc, Process{
			Pid:         p.Pid,
			Name:        "coucou",
			Interpreter: "golang",
			RestartTime: 0,
			CreatedAt:   pm2io.startTime.Unix(),
			ExecMode:    "fork_mode",
			Watching:    false,
			PmUptime:    pm2io.startTime.UnixNano(),
			Status:      "online",
			PmID:        0,
			CPU:         int(cp),
			Memory:      m.Alloc,
			NodeEnv:     "production",
			AxmActions:  pm2io.AxmActions,
			AxmOptions:  options,
		})
	}

	log.Println("nb proc", len(kmProc))

	status := PayLoad{
		Data: Data{
			Process: kmProc,
			Server: Server{
				Loadavg:     []float64{0, 0, 0},
				TotalMem:    900,
				FreeMem:     800,
				Hostname:    pm2io.Name,
				Uptime:      pm2io.startTime.Unix(),
				Pm2Version:  "10.0.0",
				Type:        "golang",
				Interaction: true,
			},
		},
		Active:     true,
		ServerName: pm2io.Name,
		InternalIP: "Coucou",
		Protected:  false,
		RevCon:     true,
	}
	msg := Message{
		Payload: status,
		Channel: "status",
	}

	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
	}
	//log.Println(string(b))

	pm2io.ws.WriteMessage(websocket.TextMessage, b)
}
