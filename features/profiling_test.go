package features_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/features"
)

func TestProfiling(t *testing.T) {
	t.Run("Should start CPU profiling", func(t *testing.T) {
		err := features.StartCPUProfile()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Should stop CPU profiling", func(t *testing.T) {
		res, err := features.StopCPUProfile()
		if err != nil {
			t.Fatal(err)
		}
		if len(res) == 0 {
			t.Fatal("result is empty")
		}
	})

	t.Run("Should heapdump", func(t *testing.T) {
		res, err := features.HeapDump()
		if err != nil {
			t.Fatal(err)
		}
		if len(res) == 0 {
			t.Fatal("result is empty")
		}
	})
}
