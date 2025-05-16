package types

import "time"

type AgentResponse struct {
	Id       string    `json:"id"`
	Hostname string    `json:"username"`
	Seen     time.Time `json:"seen"`
}
