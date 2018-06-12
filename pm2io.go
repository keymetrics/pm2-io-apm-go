package pm2io

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/keymetrics/pm2-io-apm-go/features"
	"github.com/keymetrics/pm2-io-apm-go/features/metrics"
	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/keymetrics/pm2-io-apm-go/structures"
	"github.com/shirou/gopsutil/process"
)

var version = "3.0.0-go"

type Pm2Io struct {
	Name        string
	Notifier    *features.Notifier
	transporter *services.Transporter

	startTime    time.Time
	lastCpuTotal float64
}

func (pm2io *Pm2Io) init() {
}

func (pm2io *Pm2Io) Start(publicKey string, privateKey string, name string) {
	pm2io.Name = name
	pm2io.transporter = &services.Transporter{}
	pm2io.Notifier = &features.Notifier{
		Transporter: pm2io.transporter,
	}
	services.AddMetric(metrics.GoRoutines())
	services.AddMetric(metrics.CGoCalls())
	pm2io.transporter.Connect(publicKey, privateKey, name, version)

	/*services.AddAction(&structures.Action{
		ActionName: "km:heapdump",
		ActionType: "internal",
		Callback: func() string {
			log.Println("MEMORY PROFIIIIIIIIILING")
			return ""
		},
	})
	services.AddAction(&structures.Action{
		ActionName: "km:cpuprofiling:start",
		ActionType: "internal",
		Callback: func() string {
			log.Println("CPUUUUUUUU PROFIIIIIIIIILING start")
			return ""
		},
	})
	services.AddAction(&structures.Action{
		ActionName: "km:cpuprofiling:stop",
		ActionType: "internal",
		Callback: func() string {
			log.Println("CPUUUUUUUU PROFIIIIIIIIILING stop")
			return ""
		},
	})*/

	pm2io.startTime = time.Now()

	go func() {
		ticker := time.NewTicker(time.Second)
		log.Println("created status ticker")
		for {
			<-ticker.C
			pm2io.SendStatus()
		}
	}()
}

func (pm2io *Pm2Io) SendStatus() {
	p, _ := process.NewProcess(int32(os.Getpid()))
	cp, _ := metrics.CPUPercent()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	kmProc := []structures.Process{}

	options := structures.Options{
		HeapDump:     false,
		Profiling:    false,
		CustomProbes: true,
		Error:        true,
		Errors:       true,
	}

	kmProc = append(kmProc, structures.Process{
		Pid:         p.Pid,
		Name:        pm2io.Name,
		Interpreter: "golang",
		RestartTime: 0,
		CreatedAt:   pm2io.startTime.UnixNano() / int64(time.Millisecond),
		ExecMode:    "fork_mode",
		Watching:    false,
		PmUptime:    pm2io.startTime.UnixNano() / int64(time.Millisecond),
		Status:      "online",
		PmID:        0,
		CPU:         int(cp),
		Memory:      m.Alloc,
		NodeEnv:     "production",
		AxmActions:  services.Actions,
		AxmMonitor:  services.GetMetricsAsMap(),
		AxmOptions:  options,
	})

	hostname, err := os.Hostname()
	if err != nil {
		hostname = pm2io.Name
	}
	pm2io.transporter.Send("status", structures.Status{
		Process: kmProc,
		Server: structures.Server{
			Loadavg:     metrics.CPULoad(),
			TotalMem:    metrics.TotalMem(),
			Hostname:    hostname,
			Uptime:      (time.Now().UnixNano()-pm2io.startTime.UnixNano())/int64(time.Millisecond) + 600000,
			Pm2Version:  version,
			Type:        "golang",
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
