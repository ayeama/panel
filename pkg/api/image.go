package api

type ImageResponse struct {
	Id         string            `json:"id"`
	Reference  string            `json:"reference"`
	Repository string            `json:"repository"`
	Tag        string            `json:"tag"`
	Env        map[string]string `json:"env"`
}

type ImageListResponse struct {
	Items []ImageResponse `json:"items"`
}
