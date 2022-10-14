package main

import (
	"C"
	"fmt"

	"gocv.io/x/gocv"
)

// RGBA pixel
type Pixel struct {
	R float64
	G float64
	B float64
	A float64
}

// Gray pixel
type PixelGray struct {
	pix float64
}

func main() {

	fmt.Print("\n\n#############\n")
	fmt.Println("    BEGIN")
	fmt.Print("#############\n\n\n")

	var data = make(map[string][]gocv.Mat) // name of the user : array of images
	var numUsers int
	var nameUsers []string

	// read the images
	data, numUsers, nameUsers = getData()

	var PRNUS = make(map[string][][][]PixelGray) // [layer B/G/R][rows prnu][columns prnu]

	// EXTRACTION

	for i := 0; i < numUsers; i++ {
		// does it one time per user
		// if the last parameter is "true" => the function will check the results
		PRNUS[nameUsers[i]] = extraction(data[nameUsers[i]], nameUsers[i], false)
	}

	var prnusB, prnusG, prnusR [][][]PixelGray // [user][rows][columns]

	for i := 0; i < numUsers; i++ {
		prnusUser := PRNUS[nameUsers[i]] // PRNUS B, G, R (each one is an matrix[rows][columns])
		prnusB = append(prnusB, prnusUser[0])
		prnusG = append(prnusG, prnusUser[1])
		prnusR = append(prnusR, prnusUser[2])
	}

	// ENCODED AGGREGATION

	res1 := encryption(prnusB, numUsers)
	res2 := encryption(prnusG, numUsers)
	res3 := encryption(prnusR, numUsers)

	// AGGREGATION WITHOUT ENCODING

	agreg1 := agregation(prnusB, numUsers)
	agreg2 := agregation(prnusG, numUsers)
	agreg3 := agregation(prnusR, numUsers)

	// COMPARISON BETWEEN BOTH AGGREGATIONS

	checkResults3(res1, agreg1, "B")
	checkResults3(res2, agreg2, "G")
	checkResults3(res3, agreg3, "R")

	fmt.Print("\n\n##############\n")
	fmt.Println("    FINISH")
	fmt.Print("##############\n\n")

}
