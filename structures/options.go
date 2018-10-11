package structures

// Options like AxmOptions
type Options struct {
	DefaultActions bool `json:"default_actions"`
	CustomProbes   bool `json:"custom_probes"`
	Profiling      bool `json:"profiling"`
	HeapDump       bool `json:"heapdump"`
	Apm            Apm  `json:"apm"`
}

type Apm struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}
