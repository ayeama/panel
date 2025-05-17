package types

type ServerCreateRequest struct {
	Name string `json:"name"`
}

type ServerUpdateRequest struct {
	Name string `json:"name,omitempty"`
}

type ServerResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ServerStatResponse struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}
