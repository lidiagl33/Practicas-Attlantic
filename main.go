package main

import (
	"C"

	"gocv.io/x/gocv"
)
import "fmt"

type Pixel struct {
	R float64
	G float64
	B float64
	A float64
}

type PixelGray struct {
	pix float64
}

func main() {

	var data = make(map[string][]gocv.Mat)
	var numUsers int
	var nameUsers []string

	data, numUsers, nameUsers = getData()

	var PRNUS = make(map[string][][][]PixelGray) // [layer B/G/R][rows prnu][columns prnu]

	// EXTRACTION

	for i := 0; i < numUsers; i++ {
		PRNUS[nameUsers[i]] = extraction(data[nameUsers[i]], nameUsers[i], false)
	}

	var prnusB, prnusG, prnusR [][][]PixelGray // [user][rows][columns]

	for i := 0; i < numUsers; i++ {
		prnusUser := PRNUS[nameUsers[i]] // PRNUS B, G, R (each one is an matrix[][])
		prnusB = append(prnusB, prnusUser[0])
		prnusG = append(prnusG, prnusUser[1])
		prnusR = append(prnusR, prnusUser[2])
	}

	// ENCRYPTION

	res1 := getParameters(prnusB, numUsers)
	res2 := getParameters(prnusG, numUsers)
	res3 := getParameters(prnusR, numUsers)

	// AGREGATION

	agreg1 := agregation(prnusB, numUsers)
	agreg2 := agregation(prnusG, numUsers)
	agreg3 := agregation(prnusR, numUsers)

	checkResults3(res1, agreg1, "B")
	checkResults3(res2, agreg2, "G")
	checkResults3(res3, agreg3, "R")

	fmt.Print("\nFINISH\n")

}
