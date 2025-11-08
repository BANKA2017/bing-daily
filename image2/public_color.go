//go:build !private_palette

package image2

import (
	"image"
	"image/color"
	"math"
	"strconv"

	m "github.com/ericpauley/go-quantize/quantize"
)

// By Google Gemini

type RGB struct {
	R, G, B uint8
}

type Lab struct {
	L, A, B float64
}

func RGBToLab(r, g, b uint8) Lab {
	// Step 1: normalize to [0,1]
	R := float64(r) / 255.0
	G := float64(g) / 255.0
	B := float64(b) / 255.0

	// Step 2: gamma correction
	R = gammaCorrect(R)
	G = gammaCorrect(G)
	B = gammaCorrect(B)

	// Step 3: convert to XYZ (D65)
	X := R*0.4124 + G*0.3576 + B*0.1805
	Y := R*0.2126 + G*0.7152 + B*0.0722
	Z := R*0.0193 + G*0.1192 + B*0.9505

	// Step 4: normalize to reference white D65
	X /= 0.95047
	Y /= 1.00000
	Z /= 1.08883

	// Step 5: convert XYZ to Lab
	f := func(t float64) float64 {
		if t > 0.008856 {
			return math.Cbrt(t)
		}
		return 7.787*t + 16.0/116.0
	}

	fx := f(X)
	fy := f(Y)
	fz := f(Z)

	L := 116.0*fy - 16.0
	A := 500.0 * (fx - fy)
	Bb := 200.0 * (fy - fz)

	return Lab{L, A, Bb}
}

func gammaCorrect(c float64) float64 {
	if c > 0.04045 {
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	return c / 12.92
}

func GetColor(img image.Image) ([]string, [][3]uint8, [][3]float64) {
	var colors []string
	var lab [][3]float64
	var rgb [][3]uint8

	q := m.MedianCutQuantizer{}
	p := q.Quantize(make([]color.Color, 0, 5), img)
	for _, c := range p {
		tmpColorR, tmpColorG, tmpColorB, _ := c.RGBA()
		rgb = append(rgb, [3]uint8{
			uint8(tmpColorR), uint8(tmpColorG), uint8(tmpColorB),
		})
		tmpLab := RGBToLab(uint8(tmpColorR), uint8(tmpColorG), uint8(tmpColorB))

		lab = append(lab, [3]float64{
			tmpLab.L, tmpLab.A, tmpLab.B,
		})
		colors = append(colors, padStart(strconv.FormatInt(int64(tmpColorR>>8), 16), 2, "0")+padStart(strconv.FormatInt(int64(tmpColorG>>8), 16), 2, "0")+padStart(strconv.FormatInt(int64(tmpColorB>>8), 16), 2, "0"))
	}

	return colors, rgb, lab
}
