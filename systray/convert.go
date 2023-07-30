package systray

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

func convert(b []byte) (o []byte, err error) {
	var img image.Image
	mtype := mimetype.Detect(b)
	r := bytes.NewReader(b)
	switch mtype.Extension() {
	case ".png":
		img, err = png.Decode(r)
	case ".jpg":
		img, err = jpeg.Decode(r)
	case ".gif":
		img, err = gif.Decode(r)
	case ".webp":
		img, err = webp.Decode(r)
	case ".bmp":
		img, err = bmp.Decode(r)
	case ".ico":
		return b, nil
	default:
		return b, errors.New("unsupported image format")
	}
	if err != nil {
		return b, err
	}
	w := bytes.NewBuffer([]byte(""))
	err = ico.Encode(w, img)
	if err != nil {
		return b, err
	}
	return w.Bytes(), err
}
