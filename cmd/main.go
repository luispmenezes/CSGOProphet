package main

import (
	"csgo_prophet/usecase/web_crawler"
	"log"
)

func main() {
	demoLinks := web_crawler.GetDemoLinks("2019-10-14", "2019-10-15", 1, true)

	for _, element := range demoLinks {
		log.Println(element)
	}
}
