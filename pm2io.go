package pm2io

import (
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/keymetrics/pm2-io-apm-go/features"
	"github.com/keymetrics/pm2-io-apm-go/features/metrics"
	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/keymetrics/pm2-io-apm-go/structures"
	"github.com/shirou/gopsutil/process"
)

var version = "0.0.1-go"

type Pm2Io struct {
	Config *structures.Config

	Notifier    *features.Notifier
	transporter *services.Transporter

	StatusOverrider func() *structures.Status

	serverName   string
	hostname     string
	startTime    time.Time
	lastCpuTotal float64

	cpuProfilePath string
	memProfilePath string
}

func (pm2io *Pm2Io) init() {
}

func (pm2io *Pm2Io) Start() {
	realHostname, err := os.Hostname()
	pm2io.hostname = realHostname
	serverName := ""
	if err != nil || pm2io.Config.Name != "" {
		serverName = pm2io.Config.Name
	} else {
		random, err := randomHex(5)
		if err == nil {
			serverName = realHostname + "_" + random
		} else {
			serverName = random
		}
	}

	pm2io.serverName = serverName

	node := pm2io.Config.Node
	defaultNode := "root.keymetrics.io"
	if node == nil {
		node = &defaultNode
	}

	pm2io.transporter = services.NewTransporter(pm2io.Config, version, realHostname, pm2io.serverName, *node)
	pm2io.Notifier = &features.Notifier{
		Transporter: pm2io.transporter,
	}
	metrics.InitMetricsMemStats()
	services.AttachHandler(metrics.Handler)
	services.AddMetric(metrics.GoRoutines())
	services.AddMetric(metrics.CgoCalls())
	services.AddMetric(metrics.GlobalMetricsMemStats.NumGC)
	services.AddMetric(metrics.GlobalMetricsMemStats.NumMallocs)
	services.AddMetric(metrics.GlobalMetricsMemStats.NumFree)
	services.AddMetric(metrics.GlobalMetricsMemStats.HeapAlloc)
	services.AddMetric(metrics.GlobalMetricsMemStats.Pause)

	pm2io.transporter.Connect()

	services.AddAction(&structures.Action{
		ActionName: "km:heapdump",
		ActionType: "internal",
		Callback: func(payload map[string]interface{}) string {
			f, _ := ioutil.TempFile("/tmp", "heapdump")
			pprof.WriteHeapProfile(f)
			r, _ := ioutil.ReadFile(f.Name())
			pm2io.transporter.Send("profilings", structures.NewProfilingResponse(string(r), "heapdump"))
			log.Println(payload)
			return ""
		},
	})
	services.AddAction(&structures.Action{
		ActionName: "km:cpu:profiling:start",
		ActionType: "internal",
		Callback: func(payload map[string]interface{}) string {
			f, _ := ioutil.TempFile("/tmp", "cpuprofile")
			pm2io.cpuProfilePath = f.Name()
			pprof.StartCPUProfile(f)

			if payload["opts"] != nil {
				go func() {
					timeout := payload["opts"].(map[string]interface{})["timeout"]
					time.Sleep(time.Duration(timeout.(float64)) * time.Millisecond)
					pm2io.stopAndSendCpuProfiling()
				}()
			}
			return ""
		},
	})
	services.AddAction(&structures.Action{
		ActionName: "km:cpu:profiling:stop",
		ActionType: "internal",
		Callback: func(payload map[string]interface{}) string {
			pm2io.stopAndSendCpuProfiling()
			return ""
		},
	})

	pm2io.startTime = time.Now()

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
			pm2io.SendStatus()
		}
	}()
}

func (pm2io *Pm2Io) stopAndSendCpuProfiling() {
	pprof.StopCPUProfile()
	r, _ := ioutil.ReadFile(pm2io.cpuProfilePath)
	pm2io.transporter.Send("profilings", structures.NewProfilingResponse(string(r), "cpuprofile"))
}

func (pm2io *Pm2Io) RestartTransporter() {
	pm2io.transporter.CloseAndReconnect()
}

func (pm2io *Pm2Io) SendStatus() {
	if pm2io.StatusOverrider != nil {
		pm2io.transporter.Send("status", pm2io.StatusOverrider())
		return
	}
	p, _ := process.NewProcess(int32(os.Getpid()))
	cp, _ := metrics.CPUPercent()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	kmProc := []structures.Process{}

	options := structures.Options{
		HeapDump:     true,
		Profiling:    true,
		CustomProbes: true,
		Error:        true,
		Errors:       true,
		PmxVersion:   "2.4.1",
	}

	kmProc = append(kmProc, structures.Process{
		Pid:         p.Pid,
		Name:        pm2io.Config.Name,
		Interpreter: "golang",
		RestartTime: 0,
		CreatedAt:   pm2io.startTime.UnixNano() / int64(time.Millisecond),
		ExecMode:    "fork_mode",
		Watching:    false,
		PmUptime:    pm2io.startTime.UnixNano() / int64(time.Millisecond),
		Status:      "online",
		PmID:        0,
		CPU:         cp,
		Memory:      m.Alloc,
		NodeEnv:     "production",
		AxmActions:  services.Actions,
		AxmMonitor:  services.GetMetricsAsMap(),
		AxmOptions:  options,
	})

	pm2io.transporter.Send("status", structures.Status{
		Process: kmProc,
		Server: structures.Server{
			Loadavg:     metrics.CPULoad(),
			TotalMem:    metrics.TotalMem(),
			Hostname:    pm2io.hostname,
			Uptime:      (time.Now().UnixNano()-pm2io.startTime.UnixNano())/int64(time.Millisecond) + 600000,
			Pm2Version:  version,
			Type:        runtime.GOOS,
			Interaction: true,
			CPU: structures.CPU{
				Number: runtime.NumCPU(),
				Info:   metrics.CPUName(),
			},
			NodeVersion: runtime.Version(),
		},
	})
}

func (pm2io *Pm2Io) Panic(err error) {
	pm2io.Notifier.Error(err)
	panic(err)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
