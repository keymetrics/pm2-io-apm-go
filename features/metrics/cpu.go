package metrics

import (
	"os"
	"time"

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
