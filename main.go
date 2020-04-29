package main

import (
	"etsy"
	"fmt"
	"utils"
)

func main() {
	parameters := utils.GetListingsDetails("listings.csv")
	client := etsy.NewClient()
	for _, params := range parameters {
		image := utils.GetImageFile(fmt.Sprintf("./images/%s", params.ImageName))
		client.AddListings(params, image)
	}
}
