package services

import "github.com/keymetrics/pm2-io-apm-go/structures"

var Actions []*structures.Action

// AddAction add an action to global Actions array
func AddAction(action *structures.Action) {
	Actions = append(Actions, action)
}

// CallAction with specific name (like pmx)
func CallAction(name string, payload map[string]interface{}) *string {
	for _, i := range Actions {
		if i.ActionName == name {
			response := i.Callback(payload)
			return &response
		}
	}
	return nil
}
