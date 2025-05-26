package types

type ManifestVariableRequest struct {
	Name    string   `json:"name"`
	Options []string `json:"options"`
}

type ManifestCreateRequest struct {
	Name      string `json:"name"`
	Variant   string `json:"variant"`
	Version   uint   `json:"version"`
	Variables []ManifestVariableRequest
}

type ManifestVariableResponse struct {
	Name    string   `json:"name"`
	Options []string `json:"options"`
}

type ManifestResponse struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Variant *string `json:"variant"`
	Version uint    `json:"version"`

	Variables []ManifestVariableResponse `json:"variables"`
}
