package api

type Instance struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Status    string            `json:"status"`
	Ports     map[string]string `json:"ports"`
	Resources InstanceResources `json:"resources"`
	Webhooks  []string          `json:"webhooks"`
}

type InstanceResources struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Disk   float64 `json:"disk"`
}

type InstanceCreateRequest struct {
	Image     string            `json:"image"`
	Resources InstanceResources `json:"resources"`
	Webhooks  []string          `json:"webhooks"`
}

type InstanceStat struct {
	CPUPercent     float64 `json:"cpu_percent"`
	MemoryPercent  float64 `json:"memory_percent"`
	DiskPercent    float64 `json:"disk_percent"`
	NetworkTXBytes uint64  `json:"network_tx_bytes"`
	NetworkRXBytes uint64  `json:"network_rx_bytes"`
}
