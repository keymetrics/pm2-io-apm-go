package structures

// Action like AxmAction
type Action struct {
	ActionName string        `json:"action_name"`
	Callback   func() string `json:"-"`
}
