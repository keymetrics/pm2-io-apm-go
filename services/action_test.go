package services_test

import (
	"testing"

	"github.com/keymetrics/pm2-io-apm-go/services"
	"github.com/keymetrics/pm2-io-apm-go/structures"
)

func TestAction(t *testing.T) {
	var action structures.Action

	t.Run("Create action", func(t *testing.T) {
		action = structures.Action{
			ActionName: "MyAction",
			Callback: func() string {
				return "GOOD"
			},
		}
	})

	t.Run("Add action to service", func(t *testing.T) {
		services.AddAction(&action)
	})

	t.Run("Must return correct value", func(t *testing.T) {
		resp := services.CallAction("MyAction")
		if resp == nil {
			t.Fatal("response is nil")
		}
	})

	t.Run("Must return nil for unknown action call", func(t *testing.T) {
		resp := services.CallAction("Unknown")
		if resp != nil {
			t.Fatal("response is not nil")
		}
	})
}
