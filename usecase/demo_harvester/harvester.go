package demo_harvester

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

	var roundList []model.Round

	p.RegisterEventHandler(func(e events.MatchStart) {
		log.Println("Match Start")
	})

	p.RegisterEventHandler(func(e events.RoundStart) {
		//log.Println("Round Start  ")
		var currentRound model.Round
		currentRound.StartTick = p.GameState().IngameTick()
		roundList = append(roundList, currentRound)
	})

	p.RegisterEventHandler(func(e events.RoundEndOfficial) {
		//log.Println("Round End Official")
	})

	p.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		//Calculate Economy Event
	})

	p.RegisterEventHandler(func(e events.RoundEnd) {
		if p.GameState().IsMatchStarted() {
			//log.Println("Round End")
			var currentRound = roundList[len(roundList)-1]
			currentRound.Index = p.GameState().TotalRoundsPlayed()
			currentRound.Winner = e.WinnerState.ClanName
			currentRound.EndTick = p.GameState().IngameTick()
			currentRound.EndReason = int(e.Reason)
		}
	})

	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}
}
