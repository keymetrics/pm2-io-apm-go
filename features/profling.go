package features

import (
	"io/ioutil"
	"runtime/pprof"
)

var currentPath string

// HeapDump return a binary string of a pprof
func HeapDump() (string, error) {
	f, err := ioutil.TempFile("/tmp", "heapdump")
	if err != nil {
		return "", err
	}
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		return "", err
	}
	r, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// StartCPUProfile start a CPU profiling and write it to file
func StartCPUProfile() error {
	f, err := ioutil.TempFile("/tmp", "cpuprofile")
	if err != nil {
		return err
	}
	currentPath = f.Name()
	return pprof.StartCPUProfile(f)
}

// StopCPUProfile stop the profiling and read the file
func StopCPUProfile() (string, error) {
	pprof.StopCPUProfile()
	r, err := ioutil.ReadFile(currentPath)
	if err != nil {
		return "", err
	}
	return string(r), nil
}
