package main

import (
	"C"

	"gocv.io/x/gocv"
)

func extraction(images []gocv.Mat, userName string, check bool) [][][]PixelGray {

	numImg := len(images)

	var residualsB [][][]PixelGray // [image][rows][columns]
	var residualsG [][][]PixelGray
	var residualsR [][][]PixelGray

	var layersB []gocv.Mat // [layer B of the image 0, 1, ...]
	var layersG []gocv.Mat
	var layersR []gocv.Mat

	for i := 0; i < numImg; i++ {

		imgB := gocv.NewMat()
		imgG := gocv.NewMat()
		imgR := gocv.NewMat()

		gocv.ExtractChannel(images[i], &imgB, 0)
		gocv.ExtractChannel(images[i], &imgG, 1)
		gocv.ExtractChannel(images[i], &imgR, 2)

		layersB = append(layersB, imgB)
		layersG = append(layersG, imgG)
		layersR = append(layersR, imgR)

	}

	// CALCULATE THE NUMERATOR AND DENOMINATOR WITHOUT THE ADDITION

	arrayNum1, arrayR1, originalSizes := calculatePRNU_1(layersB, &residualsB) // B
	arrayNum2, arrayR2, originalSizes := calculatePRNU_1(layersG, &residualsG) // G
	arrayNum3, arrayR3, originalSizes := calculatePRNU_1(layersR, &residualsR) // R

	maxLengthX, maxLengthY := calculateMaxLength(originalSizes)

	// ADDITION OF THE NUMERATOR AND DENOMINATOR

	var pixSumNum1 = make([][]PixelGray, maxLengthY)
	for i := 0; i < len(pixSumNum1); i++ {
		pixSumNum1[i] = make([]PixelGray, maxLengthX)
	}
	var pixSumDen1 = make([][]PixelGray, maxLengthY)
	for i := 0; i < len(pixSumDen1); i++ {
		pixSumDen1[i] = make([]PixelGray, maxLengthX)
	}

	var pixSumNum2 = make([][]PixelGray, maxLengthY)
	for i := 0; i < len(pixSumNum2); i++ {
		pixSumNum2[i] = make([]PixelGray, maxLengthX)
	}
	var pixSumDen2 = make([][]PixelGray, maxLengthY)
	for i := 0; i < len(pixSumDen2); i++ {
		pixSumDen2[i] = make([]PixelGray, maxLengthX)
	}

	var pixSumNum3 = make([][]PixelGray, maxLengthY)
	for i := 0; i < len(pixSumNum3); i++ {
		pixSumNum3[i] = make([]PixelGray, maxLengthX)
	}
	var pixSumDen3 = make([][]PixelGray, maxLengthY)
	for i := 0; i < len(pixSumDen3); i++ {
		pixSumDen3[i] = make([]PixelGray, maxLengthX)
	}

	pixSumNum1, pixSumDen1 = calculatePRNU_2(arrayNum1, arrayR1, maxLengthX, maxLengthY)
	pixSumNum2, pixSumDen2 = calculatePRNU_2(arrayNum2, arrayR2, maxLengthX, maxLengthY)
	pixSumNum3, pixSumDen3 = calculatePRNU_2(arrayNum3, arrayR3, maxLengthX, maxLengthY)

	// CALCULATE K (ESTIMATION OF THE PRNU)

	pixK1 := calculateK(pixSumNum1, pixSumDen1) // PRNU B
	pixK2 := calculateK(pixSumNum2, pixSumDen2) // PRNU G
	pixK3 := calculateK(pixSumNum3, pixSumDen3) // PRNU R

	arrayPRNUs := [][][]PixelGray{pixK1, pixK2, pixK3}

	if check {

		printResults(5, 0, userName)

		// CHECKING K (BETWEEN -1 Y 1)

		checkPRNU_1(pixK1, "PRNU B")
		checkPRNU_1(pixK2, "PRNU G")
		checkPRNU_1(pixK3, "PRNU R")

		// CHECKING K (SCALAR PRODUCT: IMAGE AND PRNU)

		printResults(0, 0, "")

		for i := 0; i < len(layersB); i++ {
			printResults(1, i, "")
			checkResults1(layersB[i], pixK1)
		}

		for i := 0; i < len(layersG); i++ {
			printResults(2, i, "")
			checkResults1(layersG[i], pixK2)
		}

		for i := 0; i < len(layersR); i++ {
			printResults(3, i, "")
			checkResults1(layersR[i], pixK3)
		}

		// CHECKING K (SCALAR PRODUCT: RESIDUAL AND PRNU)

		printResults(4, 0, "")

		for i := 0; i < len(residualsB); i++ {
			printResults(1, i, "")
			checkResults2(residualsB[i], pixK1)
		}

		for i := 0; i < len(residualsG); i++ {
			printResults(2, i, "")
			checkResults2(residualsG[i], pixK2)
		}

		for i := 0; i < len(residualsR); i++ {
			printResults(3, i, "")
			checkResults2(residualsR[i], pixK3)
		}

	}

	return arrayPRNUs
}
