package harvester

import (
	"log"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang"
	events "github.com/markus-wa/demoinfocs-golang/events"
)

func ProcessDemo(demoFilePath string) {
	f, err := os.Open(demoFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)

	// Register handler on kill events
	p.RegisterEventHandler(func(e events.Kill) {
		log.Println("KILL!!")
		log.Println(e.Killer.Name)
	})

	// Register handler on kill events
	p.RegisterEventHandler(func(e events.Kill) {
		log.Println("KILL!!")
		log.Println(e.Killer.Name)
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}
}
