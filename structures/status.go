package structures

// Status packet
type Status struct {
	Process []Process `json:"process"`
	Server  Server    `json:"server"`
}
