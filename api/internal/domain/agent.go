package domain

import "time"

type AgentStat struct {
	Uptime float64
	Cpu    float64
	Memory float64
	Time   time.Time
}

func (s *AgentStat) Online() bool {
	return time.Since(s.Time) <= (time.Second * 5)
}
