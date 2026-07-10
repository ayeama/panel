package types

type InstanceLabel string

const (
	InstanceLabelID      string = "com.github.ayeama.panel.instance.id"
	InstanceLabelWebhook string = "com.github.ayeama.panel.instance.webhook"
)

type Instance struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`

	Ports map[string]string `json:"ports"`

	Webhook string `json:"webhook"`
}

type InstanceStat struct {
	CpuPercent     float64 `json:"cpu_percent"`
	MemoryPercent  float64 `json:"memory_percent"`
	DiskPercent    float64 `json:"disk_percent"`
	NetworkTxBytes uint64  `json:"network_tx_bytes"`
	NetworkRxBytes uint64  `json:"network_rx_bytes"`
}
