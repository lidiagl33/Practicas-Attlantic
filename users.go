package main

import (
	"fmt"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

func getData() (map[string][]gocv.Mat, int, []string) {

	var data = make(map[string][]gocv.Mat) // [userName][arrayOfImages]
	var numUsers int
	var nameUsers []string

	// 5 users maximum
	for i := 0; i < 5; i++ {

		var filesNames []string
		var images []gocv.Mat

		nameDir := "user" + strconv.Itoa(i+1) // we assume the format "user1" for the folders containing the images
		nameUsers = append(nameUsers, nameDir)

		if fileExists(nameDir) {

			numUsers++

			dir, err := os.Open(nameDir)
			if err != nil {
				fmt.Println(err)
				return nil, 0, nil
			}

			files, err := dir.Readdir(0) // read all the files inside the folder
			if err != nil {
				fmt.Println(err)
				return nil, 0, nil
			}

			for j := 0; j < len(files); j++ {
				filesNames = append(filesNames, files[j].Name())
			}

			for z := 0; z < len(filesNames); z++ {
				fmt.Printf("\tLoading %q\n", nameDir+"/"+filesNames[z])
				img := gocv.IMRead(nameDir+"/"+filesNames[z], gocv.IMReadColor) // read the image
				images = append(images, img)
			}

			data[nameDir] = images

			fmt.Printf("\nImages of %q loaded\n\n", nameDir)

		}

	}

	return data, numUsers, nameUsers

}

func fileExists(rute string) bool {

	_, err := os.Stat(rute)

	if os.IsNotExist(err) {
		return false
	}

	return true
}
