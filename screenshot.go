package main

import (
	"github.com/kbinani/screenshot"
)

// calcBrightness : calculates the brightness of the screen
func calcBrightness() float64 {
	bounds := screenshot.GetDisplayBounds(0)
	bounds.Max.X = bounds.Max.X/2 + 400
	bounds.Max.Y = bounds.Max.Y/2 + 400
	bounds.Min.X = bounds.Max.X - 800
	bounds.Min.Y = bounds.Max.Y - 800

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}

	brightness := 0.0
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			brightness += float64(r) + float64(g) + float64(b)
		}
	}
	brightness /= 65535 * float64(img.Bounds().Dx()*img.Bounds().Dy()*3)
	return brightness * 100
}
