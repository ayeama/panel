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

func (s *AgentStat) Score() float64 {
	cpuWeight := 0.4
	memoryWeight := 0.6

	// NOTE invert usage to get free
	cpuScore := (1 - s.Cpu) * cpuWeight
	memoryScore := (1 - s.Memory) * memoryWeight

	return cpuScore + memoryScore
}
