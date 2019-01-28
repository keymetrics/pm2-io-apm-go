package structures

// Process data
type Process struct {
	Name   string `json:"name"`
	Server string `json:"server"`
	Rev    string `json:"rev"`
	PmID   int    `json:"pm_id"`
}
