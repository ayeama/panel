package types

type UserCreateRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}
