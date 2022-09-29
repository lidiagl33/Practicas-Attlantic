package main

import (
	"C"
	"fmt"

	"gocv.io/x/gocv"
)

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

	getParameters(prnusB, numUsers)

	// AGREGATION

	agreg1 := agregation(prnusB, numUsers)
	agreg2 := agregation(prnusG, numUsers)
	agreg3 := agregation(prnusR, numUsers)

	fmt.Println(len(agreg1), len(agreg1[0]))
	fmt.Println(len(agreg2))
	fmt.Println(len(agreg3))

}
