package model

import "fmt"

type Round struct {
	Winner        string
	Duration      int
	EndReason     string
	KillEvents    []KillEvent
	EconomyEvents []EconomyEvent
	BombEvents    []BombEvent
}

func (r Round) String() string {
	return fmt.Sprintf("[%s, %d, %s, %v, %v, %v]",
		r.Winner, r.Duration, r.EndReason, r.KillEvents, r.EconomyEvents, r.BombEvents)
}
