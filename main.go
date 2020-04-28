package main

import (
	"etsy"
)

func main() {
	client := etsy.NewClient()
	client.AddListings()
	client.FindUser("k8mkmpig")
}
