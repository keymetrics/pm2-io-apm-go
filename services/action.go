package services

import "github.com/keymetrics/pm2-io-apm-go/structures"

var Actions []*structures.Action

// INFO: must be send as map[string]AxmMonitor

func AddAction(action *structures.Action) {
	Actions = append(Actions, action)
}

func CallAction(name string, payload map[string]interface{}) *string {
	for _, i := range Actions {
		if i.ActionName == name {
			response := i.Callback(payload)
			return &response
		}
	}
	return nil
}
