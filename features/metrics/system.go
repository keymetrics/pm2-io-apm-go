package metrics

import (
	"net"
	"os"
	"time"

	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/process"
)

var lastCpuTotal float64 = 0

func CPUPercent() (float64, error) {
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

	if totalTime <= 0 {
		return 0, nil
	}

	val := (cput.Total() - lastCpuTotal) * 100
	lastCpuTotal = cput.Total()

	return val, nil
}

func CPUName() string {
	infos, err := cpu.Info()
	if err != nil {
		return ""
	}
	return infos[0].ModelName
}

func CPULoad() []float64 {
	avg, err := load.Avg()
	if err != nil {
		return []float64{}
	}
	return []float64{avg.Load1, avg.Load5, avg.Load15}
}

func TotalMem() uint64 {
	return memory.TotalMemory()
}

func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To16() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
