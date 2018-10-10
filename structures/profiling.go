package structures

import "time"

// ProfilingRequest from KM
type ProfilingRequest struct {
	Timeout int64 `json:"timeout"`
}

// ProfilingResponse to KM (data as string)
type ProfilingResponse struct {
	Data      string `json:"data"`
	At        int64  `json:"at"`
	Initiated string `json:"initiated"`
	Duration  int    `json:"duration"`
	Type      string `json:"type"`
}

// NewProfilingResponse with default values
func NewProfilingResponse(data string, element string) ProfilingResponse {
	res := ProfilingResponse{
		Data:      data,
		At:        time.Now().UnixNano() / int64(time.Millisecond),
		Initiated: "manual",
		Duration:  0,
		Type:      element,
	}
	return res
}
