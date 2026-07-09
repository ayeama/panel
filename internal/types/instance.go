package types

type Instance struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

// type InstanceStats struct {
// 	CpuPercent    string `json:"cpu_percent"`
// 	MemoryPercent string `json:"memory_percent"`
// }
