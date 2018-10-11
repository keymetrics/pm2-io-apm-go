package features

import (
	"io/ioutil"
	"runtime/pprof"
)

var currentPath string

// HeapDump return a binary string of a pprof
func HeapDump() ([]byte, error) {
	f, err := ioutil.TempFile("/tmp", "heapdump")
	if err != nil {
		return nil, err
	}
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		return nil, err
	}
	r, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, err
	}
	return r, nil
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
func StopCPUProfile() ([]byte, error) {
	pprof.StopCPUProfile()
	r, err := ioutil.ReadFile(currentPath)
	if err != nil {
		return nil, err
	}
	return r, nil
}
