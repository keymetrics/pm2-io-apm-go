package structures

// StatusProcess data
type StatusProcess struct {
	Pid         int32              `json:"pid"`
	Name        string             `json:"name"`
	Server      string             `json:"server"`
	Interpreter string             `json:"interpreter,omitempty"`
	RestartTime int                `json:"restart_time,omitempty"`
	CreatedAt   int64              `json:"created_at,omitempty"`
	ExecMode    string             `json:"exec_mode,omitempty"`
	Watching    bool               `json:"watching,omitempty"`
	PmUptime    int64              `json:"pm_uptime,omitempty"`
	Status      string             `json:"status,omitempty"`
	PmID        int                `json:"pm_id"`
	CPU         float64            `json:"cpu"`
	Rev         string             `json:"rev"`
	Memory      uint64             `json:"memory"`
	UniqueID    string             `json:"unique_id"`
	NodeEnv     string             `json:"node_env,omitempty"`
	AxmActions  []*Action          `json:"axm_actions,omitempty"`
	AxmMonitor  map[string]*Metric `json:"axm_monitor,omitempty"`
	AxmOptions  Options            `json:"axm_options,omitempty"`
}
