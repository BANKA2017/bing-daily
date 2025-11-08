package image2

import (
	"bytes"
	"image"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	Blurhash "github.com/bbrks/go-blurhash"
	"github.com/jellydator/ttlcache/v3"
)

type ImageMeta struct {
	Width    int          `json:"width"`
	Height   int          `json:"height"`
	Color    []string     `json:"color"`
	Blurhash string       `json:"blurhash"`
	Lab      [][3]float64 `json:"lab"`
	Rgb      [][3]uint8   `json:"rgb"`
}

func padStart(str string, length int, pad string) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(pad, length-len(str)) + str
}

func GetImg(imgBytes []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	return img, err
}

func GetImgMeta(img image.Image, name string) (ImageMeta, error) {
	var meta ImageMeta

	imgSize := img.Bounds().Size()

	meta.Width = imgSize.X
	meta.Height = imgSize.Y

	meta.Color, meta.Rgb, meta.Lab = GetColor(img)

	hash, err := Blurhash.Encode(4, 4, img)

	if err != nil {
		return meta, err
	}

	meta.Blurhash = hash

	ImageMetaCache.Set(name, meta, ttlcache.DefaultTTL)

	return meta, nil
}

// color

// red: #E41821
// orange: #E6551C
// yellow: #F8B80A
// green: #30A842
// blue: #00A5E5
// purple: #62288C
// pink: #EA659F
// brown: #69492A
// white: #fff
// black: #000
