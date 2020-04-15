package main

import (
	"csgo_prophet/usecase/web_crawler"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	demoLinks, err := web_crawler.GetDemoLinks("2019-10-14", "2019-10-15", 1, true)

	if err != nil {
		log.Fatal(err)
	}

	for _, element := range demoLinks {
		log.Println(element)
	}
}
