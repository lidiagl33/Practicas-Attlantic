package main

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

func calculatePRNU_1(layers []gocv.Mat, residuals *[][][]PixelGray) ([][][]PixelGray, [][][]PixelGray, []image.Point) {

	var originalImg image.Image
	var denoisedImg image.Image

	imgDenoised := gocv.NewMat()
	defer imgDenoised.Close()

	var arrayNum [][][]PixelGray // numerator [image][rows][columns]
	var arrayR [][][]PixelGray   // denominator

	var originalSizes []image.Point // store the sizes of the images

	var err error

	for i := 0; i < len(layers); i++ {

		// Y
		originalImg, err = layers[i].ToImage() // it will be coverted to image.Gray to have values between 0 and 1
		if err != nil {
			return nil, nil, nil
		}

		imgGray := image.NewGray(originalImg.Bounds())

		for y := originalImg.Bounds().Min.Y; y < originalImg.Bounds().Max.Y; y++ {
			for x := originalImg.Bounds().Min.X; x < originalImg.Bounds().Max.X; x++ {
				imgGray.Set(x, y, originalImg.At(x, y)) // Set already converts into color.Gray
			}
		}

		sizeOriginal := imgGray.Bounds().Size()
		originalSizes = append(originalSizes, sizeOriginal)

		// DENOISING

		gocv.FastNlMeansDenoising(layers[i], &imgDenoised)

		// X
		denoisedImg, err = imgDenoised.ToImage() // it will be coverted to image.Gray to have values between 0 and 1
		if err != nil {
			return nil, nil, nil
		}

		denoisedImgGray := image.NewGray(denoisedImg.Bounds())

		for y := denoisedImg.Bounds().Min.Y; y < denoisedImg.Bounds().Max.Y; y++ {
			for x := denoisedImg.Bounds().Min.X; x < denoisedImg.Bounds().Max.X; x++ {
				denoisedImgGray.Set(x, y, denoisedImg.At(x, y)) // Set already converts into color.Gray
			}
		}

		sizeDenoised := denoisedImgGray.Bounds().Size()

		// convert images into array of pixels
		var pixOri [][]PixelGray = pixelArrayGray(imgGray, sizeOriginal)         // Y
		var pixDen [][]PixelGray = pixelArrayGray(denoisedImgGray, sizeDenoised) // X

		pixRes := operateWithPixelsGray(pixOri, pixDen, "-") // W=Y-X
		*residuals = append(*residuals, pixRes)

		pixNumerador := operateWithPixelsGray(pixRes, pixDen, "*") // W*X
		pixDivisor := operateWithPixelsGray(pixDen, pixDen, "*")   // R=X*X

		arrayNum = append(arrayNum, pixNumerador)
		arrayR = append(arrayR, pixDivisor)

	}

	return arrayNum, arrayR, originalSizes

}

func calculatePRNU_2(arrayNum [][][]PixelGray, arrayR [][][]PixelGray, maxLengthX int, maxLengthY int) ([][]PixelGray, [][]PixelGray) {

	// make the summations in numerator and denominator to get the result

	var pixSumNum = make([][]PixelGray, maxLengthY)
	for i := 0; i < len(pixSumNum); i++ {
		pixSumNum[i] = make([]PixelGray, maxLengthX)
	}
	var pixSumDen = make([][]PixelGray, maxLengthY)
	for i := 0; i < len(pixSumDen); i++ {
		pixSumDen[i] = make([]PixelGray, maxLengthX)
	}

	// len(arrayNum) == len(arrayR)
	for i := 0; i < len(arrayNum); i++ {
		for j := 0; j < len(arrayNum[i]); j++ {
			for k := 0; k < len(arrayNum[i][j]); k++ {
				pixSumNum[j][k].pix += arrayNum[i][j][k].pix
			}
		}
	}

	for i := 0; i < len(arrayR); i++ {
		for j := 0; j < len(arrayR[i]); j++ {
			for k := 0; k < len(arrayR[i][j]); k++ {
				pixSumDen[j][k].pix += arrayR[i][j][k].pix
			}
		}
	}

	return pixSumNum, pixSumDen
}

func calculateK(pixNum [][]PixelGray, pixDen [][]PixelGray) [][]PixelGray {

	// make the final division to get the estimated prnu

	var K [][]PixelGray

	K = operateWithPixelsGray(pixNum, pixDen, "/")

	return K

}

func checkPRNU_1(pixK [][]PixelGray, s string) {

	var mayor int
	var menor int

	for i := 0; i < len(pixK); i++ {
		for j := 0; j < len(pixK[i]); j++ {

			if pixK[i][j].pix > 1 {
				mayor++
			}
			if pixK[i][j].pix < -1 {
				menor++
			}
		}
	}

	fmt.Println(s)
	fmt.Printf("\tnumbers > 1: %d\n", mayor)
	fmt.Printf("\tnumbers < -1: %d\n\n", menor)
}

func checkPRNU_2(img *image.Gray, pixK [][]PixelGray) {

	size := img.Bounds().Size()

	pix := pixelArrayGray(img, size)

	res := scalarProduct(pix, pixK)

	fmt.Printf("\tResult: %f\n\n", res)

}
