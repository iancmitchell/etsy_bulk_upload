package main

import (
	"etsy"
	"log"
	"os"

	"github.com/gocarina/gocsv"
)

func main() {
	paramsFile, err := os.OpenFile("listings.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln("Error Opening Params File: ", err)
	}
	defer paramsFile.Close()
	var parameters []etsy.Parameters
	if err = gocsv.Unmarshal(paramsFile, &parameters); err != nil {
		log.Fatalln("Error Parsing Params File: ", err)
	}
	log.Println(parameters)
	client := etsy.NewClient()
	for _, params := range parameters {
		client.AddListings(params)
	}
}
