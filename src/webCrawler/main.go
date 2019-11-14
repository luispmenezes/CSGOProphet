package main

import (
	"log"

	"./crawler"
)

func main() {
	demoLinks := crawler.GetDemoLinks("2019-10-14", "2019-10-15", 1, true)

	for _, element := range demoLinks {
		log.Println(element)
	}
}
