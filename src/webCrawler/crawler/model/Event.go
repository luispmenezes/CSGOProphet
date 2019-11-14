package model

import "fmt"

type Event struct {
	Title     string
	StartDate string
	EndDate   string
	PrizePool string
	Location  string
	Teams     []string
	Maps      []string
}

func (e Event) String() string {
	return fmt.Sprintf("[%s, %s, %s, %s, %s, %v, %v]",
		e.Title, e.StartDate, e.EndDate, e.PrizePool, e.Location, e.Teams, e.Maps)
}
