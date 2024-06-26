package main

import (
	"image"
)

func pixelArrayGray(img *image.Gray, size image.Point) [][]PixelGray {

	// convert the image into an array with pixel gray values [0,...,1]

	var pixels = make([][]PixelGray, size.Y)

	for i := 0; i < len(pixels); i++ {
		pixels[i] = make([]PixelGray, size.X)
	}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			pixels[y][x].pix = float64(img.GrayAt(x, y).Y)
		}
	}

	return pixels

}

func operateWithPixelsGray(pix1 [][]PixelGray, pix2 [][]PixelGray, check string) [][]PixelGray {

	// does the operation indicated by "check"

	var pixResult [][]PixelGray

	if len(pix1) >= len(pix2) {
		pixResult = make([][]PixelGray, len(pix1))
	} else {
		pixResult = make([][]PixelGray, len(pix2))
	}

	if len(pix1[0]) >= len(pix2[0]) {
		for i := 0; i < len(pix1); i++ {
			pixResult[i] = make([]PixelGray, len(pix1[0]))
		}
	} else {
		for i := 0; i < len(pix2); i++ {
			pixResult[i] = make([]PixelGray, len(pix2[0]))
		}
	}

	for i := 0; i < len(pixResult); i++ {
		for j := 0; j < len(pixResult[0]); j++ {

			if check == "+" {

				y := pix1[i][j].pix + pix2[i][j].pix

				pixResult[i][j].pix = y

			} else if check == "-" {

				y := pix1[i][j].pix - pix2[i][j].pix

				pixResult[i][j].pix = y

			} else if check == "*" {

				y := pix1[i][j].pix * pix2[i][j].pix

				pixResult[i][j].pix = y

			} else if check == "/" {

				var y float64

				if pix2[i][j].pix == 0 {
					y = pix1[i][j].pix / (pix2[i][j].pix + 0.000001) // make the denominator not zero
				} else {
					y = pix1[i][j].pix / pix2[i][j].pix
				}

				pixResult[i][j].pix = y
			}

		}
	}

	return pixResult
}
