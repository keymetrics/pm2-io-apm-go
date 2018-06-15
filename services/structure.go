package services

import "github.com/keymetrics/pm2-io-apm-go/structures"

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

type Verify struct {
	PublicId  string     `json:"public_id"`
	PrivateId string     `json:"private_id"`
	Data      VerifyData `json:"data"`
}

type VerifyData struct {
	MachineName string `json:"MACHINE_NAME"`
	Cpus        int    `json:"CPUS"`   //nb thread
	Memory      uint64 `json:"MEMORY"` //bytes
	Pm2Version  string `json:"PM2_VERSION"`
	Hostname    string `json:"HOSTNAME"`
}

type VerifyResponse struct {
	Endpoints Endpoints `json:"endpoints"`
}

type Endpoints struct {
	WS string `json:"ws"`
}
