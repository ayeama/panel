package types

type ServerCreateRequest struct {
	Image string `json:"image"`
}

type ServerResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Address string `json:"address"`
}
