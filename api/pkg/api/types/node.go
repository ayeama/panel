package types

type NodeCreateRequest struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type NodeResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Uri  string `json:"uri"`
}
