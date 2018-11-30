package structures_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/structures"
)

func TestConfig(t *testing.T) {
	var config structures.Config

	t.Run("Should create a config object", func(t *testing.T) {
		config = structures.Config{
			Name: "golang_tests",
		}
	})

	t.Run("Shouldn't have a serverName", func(t *testing.T) {
		if config.ServerName != "" {
			t.Fatal("Already have a serverName")
		}
	})

	t.Run("Should set a serverName", func(t *testing.T) {
		config.InitNames()
		if config.ServerName == "" {
			t.Fatal("No serverName generated")
		}
		if config.Hostname == "" {
			t.Fatal("No hostname")
		}
	})
}
