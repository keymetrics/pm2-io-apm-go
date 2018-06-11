package structures

type Options struct {
	DefaultActions bool `json:"default_actions"`
	CustomProbes   bool `json:"custom_probes"`
	Error          bool `json:"error"`
	Errors         bool `json:"errors"`
	Profiling      bool `json:"profiling"`
	HeapDump       bool `json:"heapdump"`
}
