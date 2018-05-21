package pm2_io_apm_go

type Message struct {
	Payload PayLoad `json:"payload"`
	Channel string  `json:"channel"`
}

type MessageMap struct {
	Payload map[string]interface{} `json:"payload"`
	Channel string                 `json:"channel"`
}

// Sending
type PayLoad struct {
	Data       Data   `json:"data"`
	Active     bool   `json:"active"`
	ServerName string `json:"server_name"`
	InternalIP string `json:"internal_ip"`
	Protected  bool   `json:"protected"`
	RevCon     bool   `json:"rev_con"`
}
type Data struct {
	Process []Process `json:"process"`
	Server  Server    `json:"server"`
}
type Process struct {
	Pid         int32  `json:"pid"`
	Name        string `json:"name"`
	Interpreter string `json:"interpreter,omitempty"`
	RestartTime int    `json:"restart_time,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	ExecMode    string `json:"exec_mode"`
	Watching    bool   `json:"watching,omitempty"`
	PmUptime    int64  `json:"pm_uptime,omitempty"`
	Status      string `json:"status"`
	PmID        int    `json:"pm_id"`
	CPU         int    `json:"cpu"`
	Memory      uint64 `json:"memory"`
	Versioning  struct {
		Comment              string `json:"comment"`
		URL                  string `json:"url"`
		Revision             string `json:"revision"`
		Branch               string `json:"branch"`
		Type                 string `json:"type"`
		BranchExistsOnRemote bool   `json:"branch_exists_on_remote"`
	} `json:"versioning,omitempty"`
	NodeEnv    string      `json:"node_env,omitempty"`
	AxmActions []AxmAction `json:"axm_actions"`
	AxmMonitor struct {
	} `json:"axm_monitor,omitempty"`
	AxmOptions AxmOptions `json:"axm_options"`
}
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
type AxmOptions struct {
	DefaultActions bool `json:"default_actions"`
	CustomProbes   bool `json:"custom_probes"`
	Error          bool `json:"error"`
	Errors         bool `json:"errors"`
	Profiling      bool `json:"profiling"`
	HeapDump       bool `json:"heapdump"`
}
type CPU struct {
	Number int    `json:"number"`
	Info   string `json:"info"`
}
type AxmAction struct {
	ActionName string        `json:"action_name"`
	Callback   func() string `json:"-"`
}

// Receiving
type MessageResponse struct {
	Payload AxmActionResponse `json:"payload"`
	Channel string            `json:"channel"`
}
type AxmActionResponse struct {
	ActionName string `json:"action_name"`
	ProcessId  int    `json:"process_id"`
}
type AxmActionSucess struct {
	Success    bool   `json:"success"`
	Id         int    `json:"id"`
	ActionName string `json:"action_name"`
	Return     string `json:"return"`
}
