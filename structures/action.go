package structures

// Action like AxmAction
type Action struct {
	ActionName string        `json:"action_name"`
	ActionType string        `json:"action_type"` // default: "custom" else "internal" (like profiling, restart)
	Callback   func() string `json:"-"`
}
