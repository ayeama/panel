package types

type ImageResponse struct {
	Image string `json:"image"`

	Variables map[string]string `json:"variables"`
}
