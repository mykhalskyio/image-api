package entity

// image entity
type Image struct {
	Id                       int
	ImagePathQualityOriginal string
	ImagePathQuality75       string
	ImagePathQuality50       string
	ImagePathQuality25       string
}

type ImageUpload struct {
	Image     []byte `json:"image"`
	ImageName string `json:"image-name"`
}
