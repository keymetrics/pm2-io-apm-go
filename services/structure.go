package services

import "github.com/f-hj/pm2-io-apm-go/structures"

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
