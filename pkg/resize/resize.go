package resize

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

// resize image with n quality from the original
func ImgResize(path string, quality float64) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, filetype, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	newImage := resize.Resize(
		uint(float64(img.Bounds().Dx())*(quality/100)),
		uint(float64(img.Bounds().Dy())*(quality/100)),
		img, resize.Lanczos3)

	buf := bytes.Buffer{}

	switch filetype {
	case "png":
		err = png.Encode(&buf, newImage)
	case "gif":
		err = gif.Encode(&buf, newImage, nil)
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, newImage, nil)
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
