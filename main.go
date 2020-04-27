package main

import "etsy"

func main() {
	client := etsy.NewClient()
	client.AddListings()
}
