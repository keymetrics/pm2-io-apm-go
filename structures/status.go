package structures

// Status packet
type Status struct {
	Process []StatusProcess `json:"process"`
	Server  Server          `json:"server"`
}
