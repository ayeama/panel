package domain

// agents

type EventAgentStat struct {
	Hostname string  `json:"hostname"`
	Uptime   float64 `json:"uptime"`
	Cpu      float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
}

// servers

type EventServerCreate struct {
	Id string `json:"id"`
}

type EventServerCreated struct{}

type EventServerStart struct {
	Id          string `json:"id"`
	ContainerId string `json:"container_id"`
}

type EventServerStarted struct{}

type EventServerStop struct {
	Id          string `json:"id"`
	ContainerId string `json:"container_id"`
}

type EventServerStopped struct{}

type EventServerDelete struct {
	Id          string `json:"id"`
	ContainerId string `json:"container_id"`
}

type EventServerDeleted struct{}
