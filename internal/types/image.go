package types

type ImageLabel string

const (
	ImageLabelID          string = "com.github.ayeama.panel.image.id"
	ImageLabelVersion     string = "com.github.ayeama.panel.image.version"
	ImageLabelName        string = "com.github.ayeama.panel.image.name"
	ImageLabelDescription string = "com.github.ayeama.panel.image.description"
)

type Image struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
