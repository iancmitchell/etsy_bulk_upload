package main

import (
	"etsy"
	"log"
)

func main() {
	client := etsy.NewClient()
	res := client.AddListings()
	log.Println(res)
}
