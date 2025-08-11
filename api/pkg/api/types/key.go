package types

type KeyCreateRequest struct {
	Comment   string `json:"comment"`
	PublicKey string `json:"key"`
}

type KeyResponse struct {
	Id              string `json:"id"`
	Comment         string `json:"comment"`
	HashedPublicKey string `json:"key"`
}
