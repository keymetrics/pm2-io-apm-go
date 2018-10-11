package structures

import (
	b64 "encoding/base64"
	"time"
)

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
	Encoding  string `json:"encoding"`
}

// NewProfilingResponse with default values
func NewProfilingResponse(data []byte, element string) ProfilingResponse {
	res := ProfilingResponse{
		Data:      b64.StdEncoding.EncodeToString(data),
		At:        time.Now().UnixNano() / int64(time.Millisecond),
		Initiated: "manual",
		Duration:  0,
		Type:      element,
		Encoding:  "base64",
	}
	return res
}
