package main

import "log"

func main() {
	demoLinks := web_crawler.GetDemoLinks("2019-10-14", "2019-10-15", 1, true)

	for _, element := range demoLinks {
		log.Println(element)
	}
}
