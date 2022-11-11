package main

import (
	"C"
	"fmt"

	"math"
	"time"

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

func runTimed(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
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

	//mt.Printf("\nTime consumed to do EXTRACTION: %s\n\n", durationExtraction)

	var prnusB, prnusG, prnusR [][][]PixelGray // [user][rows][columns]

	for i := 0; i < numUsers; i++ {
		prnusUser := PRNUS[nameUsers[i]] // PRNUS B, G, R (each one is an matrix[rows][columns])
		prnusB = append(prnusB, prnusUser[0])
		prnusG = append(prnusG, prnusUser[1])
		prnusR = append(prnusR, prnusUser[2])
	}

	// ENCODED AGGREGATION

	var res1, res2, res3 [][]float64

	res1 = encryption(prnusB, numUsers)
	res2 = encryption(prnusG, numUsers)
	res3 = encryption(prnusR, numUsers)

	// AGGREGATION WITHOUT ENCODING

	var agreg1, agreg2, agreg3 [][]PixelGray

	agreg1 = agregation(prnusB, numUsers)
	agreg2 = agregation(prnusG, numUsers)
	agreg3 = agregation(prnusR, numUsers)

	// estimated prnu average (without ecryption) -> relative error
	var averageK1, averageK2, averageK3 float64

	for i := 0; i < len(res1); i++ {
		for j := 0; j < len(res1[0]); j++ {
			averageK1 += math.Abs(res1[i][j])
			averageK2 += math.Abs(res2[i][j])
			averageK3 += math.Abs(res3[i][j])
		}
	}

	averageK1 = averageK1 / float64((len(res1) * len(res1[0])))
	averageK2 = averageK2 / float64((len(res2) * len(res2[0])))
	averageK3 = averageK3 / float64((len(res3) * len(res3[0])))

	// COMPARISON BETWEEN BOTH AGGREGATIONS

	checkResults3(res1, agreg1, averageK1, "B")
	checkResults3(res2, agreg2, averageK2, "G")
	checkResults3(res3, agreg3, averageK3, "R")

	fmt.Print("\n\n##############\n")
	fmt.Println("    FINISH")
	fmt.Print("##############\n\n")

}
