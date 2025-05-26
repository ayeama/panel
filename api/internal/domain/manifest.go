package domain

type ManifestVariable struct {
	Name    string   `json:"name"`
	Options []string `json:"version"`
}

type ManifestPayload struct {
	Variables []ManifestVariable `json:"variables"`
}

type Manifest struct {
	Id      string
	Name    string
	Variant *string
	Version uint

	Payload ManifestPayload
}
