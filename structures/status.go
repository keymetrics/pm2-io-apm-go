package structures

type Status struct {
	Process []Process `json:"process"`
	Server  Server    `json:"server"`
}
