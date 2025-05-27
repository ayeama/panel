package domain

type ManifestVariable struct {
	Name    string   `json:"name"`
	Options []string `json:"version"`
}

type ManifestPayload struct {
	Image     string             `json:"image"`
	Variables []ManifestVariable `json:"variables"`
}

type Manifest struct {
	Id      string
	Name    string
	Variant *string
	Version uint

	Payload ManifestPayload
}
