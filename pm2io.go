package pm2_io_apm_go

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/f-hj/pm2-io-apm-go/features/metrics"
	"github.com/f-hj/pm2-io-apm-go/services"
	"github.com/f-hj/pm2-io-apm-go/structures"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/process"
)

var version = "3.0.0-go"

type Pm2Io struct {
	Name        string
	transporter *services.Transporter

	startTime    time.Time
	lastCpuTotal float64
}

func (pm2io *Pm2Io) init() {
}

func (pm2io *Pm2Io) Start(publicKey string, privateKey string, name string) {
	pm2io.Name = name
	pm2io.transporter = &services.Transporter{}
	pm2io.transporter.Connect(publicKey, privateKey, name, version)

	services.AddAction(&structures.Action{
		ActionName: "km:heapdump",
		Callback: func() string {
			log.Println("MEMORY PROFIIIIIIIIILING")
			return ""
		},
	})
	services.AddAction(&structures.Action{
		ActionName: "km:cpuprofiling:start",
		Callback: func() string {
			log.Println("CPUUUUUUUU PROFIIIIIIIIILING start")
			return ""
		},
	})
	services.AddAction(&structures.Action{
		ActionName: "km:cpuprofiling:stop",
		Callback: func() string {
			log.Println("CPUUUUUUUU PROFIIIIIIIIILING stop")
			return ""
		},
	})

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
		HeapDump:     true,
		Profiling:    true,
		CustomProbes: true,
		Error:        true,
		Errors:       true,
	}

	kmProc = append(kmProc, structures.Process{
		Pid:         p.Pid,
		Name:        pm2io.Name,
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
		AxmActions:  services.Actions,
		AxmMonitor:  services.GetMetricsAsMap(),
		AxmOptions:  options,
	})

	pm2io.transporter.Send("status", structures.Status{
		Process: kmProc,
		Server: structures.Server{
			Loadavg:     []float64{0, 0, 0},
			TotalMem:    900000000,
			FreeMem:     800,
			Hostname:    pm2io.Name,
			Uptime:      pm2io.startTime.Unix(),
			Pm2Version:  version,
			Type:        "golang",
			Interaction: true,
		},
	})
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (pm2io *Pm2Io) NotifyError(err error) {
	log.Println("ERRRRROOOOOOORRRR")

	stack := ""
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			stack += fmt.Sprintf("%+v", f)
		}
	} else {
		stack = fmt.Sprintf("%+v", err)
	}
	pm2io.transporter.Send("process:exception", Error{
		Message: err.Error(),
		Stack:   stack,
	})
}
