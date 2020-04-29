package utils

import (
	"etsy"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

//GetListingsDetails imports the parameters csv file.
func GetListingsDetails(filePath string) []etsy.Parameters {
	listingsFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln("Error Opening Params File: ", err)
	}
	defer listingsFile.Close()
	var listings []etsy.Parameters
	if err = gocsv.Unmarshal(listingsFile, &listings); err != nil {
		log.Fatalln("Error Parsing Params File: ", err)
	}
	return listings
}

//GetImageFile opens and image given a path to the file.
func GetImageFile(imagePath string) []byte {
	imageFile, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Error Opening Image: %s | Error: ", imagePath, err)
	}
	return imageFile
}
