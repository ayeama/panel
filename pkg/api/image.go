package api

type ImageResponse struct {
	Id         string            `json:"id"`
	Repository string            `json:"repository"`
	Tag        string            `json:"tag"`
	Env        map[string]string `json:"env"`
}

type ImageListResponse struct {
	Items []ImageResponse `json:"items"`
}
