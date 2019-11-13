package main

import (
	"log"

	"./crawler"
)

func main() {
	demoLinks := crawler.GetDemoLinks("", "", 0, true)

	for _, element := range demoLinks {
		log.Println(element)
	}
}
