package demo

import "fmt"

type Round struct {
	Index         int
	Winner        string
	StartTick     int
	EndTick       int
	EndReason     int
	KillEvents    []KillEvent
	EconomyEvents []EconomyEvent
	BombEvents    []BombEvent
}

func (r Round) String() string {
	return fmt.Sprintf("[%d %s, %d, %d, %d, %v, %v, %v]",
		r.Index, r.Winner, r.StartTick, r.EndTick, r.EndReason, r.KillEvents, r.EconomyEvents, r.BombEvents)
}
