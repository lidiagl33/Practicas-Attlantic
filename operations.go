package main

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"
)

func scalarProduct(pix, k [][]PixelGray) [][]PixelGray {

	//var result float64

	result := operateWithPixelsGray(pix, k, "*")

	/*for i := 0; i < len(pixMult); i++ {

		result += pixMult[i].pix

	}*/

	return result
}

func calculateMaxLength(originalSizes []image.Point) (int, int) {

	// maximum lenght of the images among all of them

	var maxLenghtX int
	var maxLenghtY int

	for i := 0; i < len(originalSizes); i++ {
		if i == 0 {
			maxLenghtX = originalSizes[i].X
			maxLenghtY = originalSizes[i].Y
		} else {
			if originalSizes[i].X > maxLenghtX {
				maxLenghtX = originalSizes[i].X
			}
			if originalSizes[i].Y > maxLenghtY {
				maxLenghtY = originalSizes[i].Y
			}
		}
	}

	return maxLenghtX, maxLenghtY
}

func convertToGray(img image.Image) *image.Gray {

	imgGray := image.NewGray(img.Bounds())

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			imgGray.Set(x, y, img.At(x, y))
		}
	}

	return imgGray

}

func checkResults1(img gocv.Mat, pixK [][]PixelGray) {

	imgIm, err := img.ToImage()
	if err != nil {
		return
	}

	imgImGray := convertToGray(imgIm)

	checkPRNU_2(imgImGray, pixK)

}

func checkResults2(residual [][]PixelGray, pixK [][]PixelGray) {

	res := scalarProduct(residual, pixK)

	fmt.Printf("result: %f\n\n", res)

}

func checkResults3(res [][]float64, agreg [][]PixelGray, layer string) {

	// res => Obtained result by encryption
	// agreg => Expected result without enctyption

	errorEncrypted := make([][]float64, len(res)) // len(res) == len(agreg)
	for i := range res {
		errorEncrypted[i] = make([]float64, len(res[0]))
	}

	for i := 0; i < len(res); i++ {
		for j := 0; j < len(res[0]); j++ {
			errorEncrypted[i][j] = agreg[i][j].pix - res[i][j]
		}
	}

	var averageError float64

	// calculate de average error: add up all values / number of values
	for i := 0; i < len(errorEncrypted); i++ {
		for j := 0; j < len(errorEncrypted[0]); j++ {
			averageError += errorEncrypted[i][j]
		}
	}

	averageError = averageError / float64((len(errorEncrypted) * len(errorEncrypted[0])))

	fmt.Printf("ERROR from layer %s OBTAINED = %f\n", layer, averageError)

}

func printResults(c, index int, userName string) {

	if c == 0 {
		fmt.Print("\n")
		fmt.Println("--- SCALAR PRODUCT WITH PRNUs ---")
		fmt.Print("\n")
	} else if c == 1 {
		fmt.Printf("--- checking prnu B%d ---", index+1)
		fmt.Print("\n")
	} else if c == 2 {
		fmt.Printf("--- checking prnu G%d ---", index+1)
		fmt.Print("\n")
	} else if c == 3 {
		fmt.Printf("--- checking prnu R%d ---", index+1)
		fmt.Print("\n")
	} else if c == 4 {
		fmt.Print("\n")
		fmt.Println("--- SCALAR PRODUCT WITH RESIDUALS ---")
		fmt.Print("\n")
	} else if c == 5 {
		fmt.Printf("----------------------------------------- USER: %q -------------------------------------------\n\n", userName)
	}

}
