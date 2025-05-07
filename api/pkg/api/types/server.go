package types

type ServerCreateRequest struct {
	Name string `json:"name"`
}

type ServerResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
