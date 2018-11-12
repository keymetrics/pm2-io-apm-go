package services

import "github.com/keymetrics/pm2-io-apm-go/structures"

// PayLoad is structure for receiving json data
type PayLoad struct {
	At         int64              `json:"at"`
	Data       interface{}        `json:"data"`
	Process    structures.Process `json:"process,omitempty"`
	Active     bool               `json:"active"`
	ServerName string             `json:"server_name"`
	InternalIP string             `json:"internal_ip"`
	Protected  bool               `json:"protected"`
	RevCon     bool               `json:"rev_con"`
}

// Verify is the object send to api.cloud.pm2.io to get a node to connect
type Verify struct {
	PublicId  string     `json:"public_id"`
	PrivateId string     `json:"private_id"`
	Data      VerifyData `json:"data"`
}

// VerifyData is a part of Verify
type VerifyData struct {
	MachineName string `json:"MACHINE_NAME"`
	Cpus        int    `json:"CPUS"`   //nb thread
	Memory      uint64 `json:"MEMORY"` //bytes
	Pm2Version  string `json:"PM2_VERSION"`
	Hostname    string `json:"HOSTNAME"`
}

// VerifyResponse is the object sent by api.cloud.pm2.io
type VerifyResponse struct {
	Endpoints Endpoints `json:"endpoints"`
}

// Endpoints list of api.cloud.pm2.io
type Endpoints struct {
	WS string `json:"ws"`
}
