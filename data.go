package pm2_io_apm_go

// Sending

type Error struct {
	Message string `json:"message"`
	Stack   string `json:"stack"`
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
