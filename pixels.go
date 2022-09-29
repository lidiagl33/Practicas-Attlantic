package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func pixelArray(size image.Point, img image.Image) [][]Pixel {

	var pixels [][]Pixel

	for x := 0; x < size.X; x++ {
		var col []Pixel
		for y := 0; y < size.Y; y++ {
			col = append(col, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, col)
	}

	return pixels
}

func pixelArrayGray(img *image.Gray, size image.Point) [][]PixelGray {

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

func operateWithPixels(pix1 [][]Pixel, pix2 [][]Pixel, check string) [][]Pixel {

	var pixResult = make([][]Pixel, len(pix1))

	for i := 0; i < len(pixResult); i++ {
		pixResult[i] = make([]Pixel, len(pix1[0]))
	}

	// suppose that pix1 and pix2 have the same lenght (otherwise fill with 0's ??)
	for i := 0; i < len(pix1); i++ {
		for j := 0; j < len(pix1[0]); j++ {

			if check == "+" {

				r := pix1[i][j].R + pix2[i][j].R
				g := pix1[i][j].G + pix2[i][j].G
				b := pix1[i][j].B + pix2[i][j].B
				a := pix1[i][j].A + pix2[i][j].A

				pixResult[i][j].R = r
				pixResult[i][j].G = g
				pixResult[i][j].B = b
				pixResult[i][j].A = a

			} else if check == "-" {

				r := pix1[i][j].R - pix2[i][j].R
				g := pix1[i][j].G - pix2[i][j].G
				b := pix1[i][j].B - pix2[i][j].B
				a := pix1[i][j].A - pix2[i][j].A

				pixResult[i][j].R = r
				pixResult[i][j].G = g
				pixResult[i][j].B = b
				pixResult[i][j].A = a

			} else if check == "*" {

				r := pix1[i][j].R * pix2[i][j].R
				g := pix1[i][j].G * pix2[i][j].G
				b := pix1[i][j].B * pix2[i][j].B
				a := pix1[i][j].A * pix2[i][j].A

				pixResult[i][j].R = r
				pixResult[i][j].G = g
				pixResult[i][j].B = b
				pixResult[i][j].A = a

			} else if check == "/" {

				var r, g, b, a float64

				if pix2[i][j].R == 0 {
					r = pix1[i][j].R / (pix2[i][j].R + 0.000001)
				} else {
					r = pix1[i][j].R / pix2[i][j].R
				}

				if pix2[i][j].G == 0 {
					g = pix1[i][j].G / (pix2[i][j].G + 0.000001)
				} else {
					g = pix1[i][j].G / pix2[i][j].G
				}

				if pix2[i][j].B == 0 {
					b = pix1[i][j].B / (pix2[i][j].B + 0.000001)
				} else {
					b = pix1[i][j].B / pix2[i][j].B
				}

				if pix2[i][j].A == 0 {
					a = pix1[i][j].A / (pix2[i][j].A + 0.000001)
				} else {
					a = pix1[i][j].A / pix2[i][j].A
				}

				pixResult[i][j].R = r
				pixResult[i][j].G = g
				pixResult[i][j].B = b
				pixResult[i][j].A = a
			}
		}
	}

	return pixResult
}

func operateWithPixelsGray(pix1 [][]PixelGray, pix2 [][]PixelGray, check string) [][]PixelGray {

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
					y = pix1[i][j].pix / (pix2[i][j].pix + 0.000001)
				} else {
					y = pix1[i][j].pix / pix2[i][j].pix
				}

				pixResult[i][j].pix = y
			}

		}
	}

	return pixResult
}

/*func grayToFloat(y uint8) PixelGray {
	return PixelGray{float64(y)}
}*/

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{float64(r / 257), float64(g / 257), float64(b / 257), float64(a / 257)}
}

func createImage(pixels [][]Pixel, check string) (error, image.Image) {

	rect := image.Rect(0, 0, len(pixels), len(pixels[0]))
	newImg := image.NewRGBA(rect)

	//newImg.Pix[0]

	for x := 0; x < len(pixels); x++ {
		for y := 0; y < len(pixels[0]); y++ {

			p := pixels[x][y]

			pixRgba, ok := pixelToRgba(p)
			var pixResult color.RGBA

			if check == "r" {
				pixResult = color.RGBA{pixRgba.R, 0, 0, 255}
			} else if check == "g" {
				pixResult = color.RGBA{0, pixRgba.G, 0, 255}
			} else if check == "b" {
				pixResult = color.RGBA{0, 0, pixRgba.B, 255}
			} else {
				pixResult = color.RGBA{pixRgba.R, pixRgba.G, pixRgba.B, pixRgba.A}
			}

			if ok {
				newImg.Set(x, y, pixResult)
			}
		}
	}

	name := check + ".jpg"
	outF, err := os.Create(name)

	if err != nil {
		fmt.Printf("Error creating file: %s", err)
		return err, nil
	}

	jpeg.Encode(outF, newImg, nil)

	outF.Close()

	return nil, newImg
}

func pixelToRgba(pixel Pixel) (color.RGBA, bool) {

	r := uint8(pixel.R)
	g := uint8(pixel.G)
	b := uint8(pixel.B)
	a := uint8(pixel.A)

	return color.RGBA{r, g, b, a}, true
}

func makeFunction() {

}
