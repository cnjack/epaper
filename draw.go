package epaper

import (
	"fmt"
	"image"
)

func GetBuffer(image image.Image) []byte {
	lineWidth := w / 8
	if w%8 != 0 {
		lineWidth++
	}

	size := (lineWidth * h)
	data := make([]byte, size)
	for i := 0; i < len(data); i++ {
		data[i] = 0xFF
	}
	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	if imageWidth == w && imageHeight == h {
		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				if isBlack(image, x, y) {
					pos := imageWidth - x
					data[pos/8+y*lineWidth] &= ^(0x80 >> (pos % 8))
				}
			}
		}
		return data
	}

	if imageWidth == h && imageHeight == w {
		for y := 0; y < imageHeight; y++ {
			for x := 0; x < imageWidth; x++ {
				if isBlack(image, x, y) {
					posx := y
					posy := imageWidth - (h - x - 1) - 1
					data[posx/8+posy*lineWidth] &= ^(0x80 >> (y % 8))
				}
			}
		}
		return data
	}
	fmt.Printf("Can't convert image expected %d %d but having %d %d", lineWidth, h, imageWidth, imageHeight)
	return data
}

func isBlack(image image.Image, x, y int) bool {
	r, g, b, a := getRGBA(image, x, y)
	offset := 10
	return r < 255-offset && g < 255-offset && b < 255-offset && a > offset
}

func getRGBA(image image.Image, x, y int) (int, int, int, int) {
	r, g, b, a := image.At(x, y).RGBA()
	r = r / 257
	g = g / 257
	b = b / 257
	a = a / 257

	return int(r), int(g), int(b), int(a)
}
