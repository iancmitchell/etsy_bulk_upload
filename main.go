package main

import (
	"etsy"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

//getListingsDetails imports the parameters csv file.
func getListingsDetails(filePath string) []etsy.Parameters {
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

func main() {
	parameters := getListingsDetails("listings.csv")
	client := etsy.NewClient()
	for _, params := range parameters {
		client.AddListings(params)
	}
}
