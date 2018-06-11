package structures

type Server struct {
	Loadavg     []float64 `json:"loadavg"`
	TotalMem    int64     `json:"total_mem"`
	FreeMem     int64     `json:"free_mem"`
	CPU         CPU       `json:"cpu"`
	Hostname    string    `json:"hostname"`
	Uptime      int64     `json:"uptime"`
	Type        string    `json:"type"`
	Platform    string    `json:"platform"`
	Arch        string    `json:"arch"`
	User        string    `json:"user"`
	Interaction bool      `json:"interaction"`
	Pm2Version  string    `json:"pm2_version"`
	NodeVersion string    `json:"node_version"`
	RemoteIP    string    `json:"remote_ip"`
	RemotePort  int       `json:"remote_port"`
}

type CPU struct {
	Number int    `json:"number"`
	Info   string `json:"info"`
}
